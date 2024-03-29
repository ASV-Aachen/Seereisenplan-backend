package setup

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	gocloak "github.com/ASV-Aachen/Seereisenplan-backend/modules/gocloak"
	"github.com/gofiber/fiber/v2"
)

// keycloak
var hostname string = os.Getenv("keycloak_hostname")
var clientID string = os.Getenv("keycloak_clientID")
var clientSecret string = os.Getenv("keycloak_clientSecret")
var realm string = os.Getenv("keycloak_realm")

func Get_AdminToken() (string, error) {

	url := "https://" + hostname + "/sso/auth/realms/" + realm + "/protocol/openid-connect/token"

	payload := strings.NewReader("client_id=" + clientID + "&grant_type=client_credentials&client_secret=" + clientSecret)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	if res.Status != "200 OK" {
		log.Default().Printf("client_id=" + clientID + "&grant_type=client_credentials&client_secret=" + clientSecret)
		return "", errors.New("")
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	answer := gocloak.AdminToken{}
	json.Unmarshal(body, &answer)

	log.Default().Printf(res.Status)

	return answer.AccessToken, nil
}

func Get_UserID(token string) (gocloak.UserInfo, error) {
	path := "https://" + hostname + "/sso/auth/realms/" + realm + "/protocol/openid-connect/userinfo"

	// payload := strings.NewReader("client_id=backend-check&grant_type=client_credentials")

	req, _ := http.NewRequest("GET", path, nil)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, _ := http.DefaultClient.Do(req)

	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		log.Default().Printf(resp.Status)
		log.Default().Printf(path)
		log.Default().Printf("[" + token + "]")

		return gocloak.UserInfo{}, errors.New("unathorized")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var answer gocloak.UserInfo
	json.Unmarshal(body, &answer)

	return answer, nil
}

func Get_UserGroups(token string, ID string) (gocloak.GroupToken, error) {

	url := "https://" + hostname + "/sso/auth/admin/realms/" + realm + "/users/" + ID + "/groups"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	if res.Status != "200 OK" {
		log.Default().Printf(res.Status)
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	answer := gocloak.GroupToken{}
	json.Unmarshal(body, &answer)

	return answer, nil
}

func Check_IsUserPartOfGroup(gruppen []string, UserGruppen gocloak.GroupToken) bool {
	for _, groupName := range gruppen {
		for _, token := range UserGruppen {
			if groupName == token.Name {
				return true
			}
		}
	}
	return false
}

func Check_IsUserTakel(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		log.Fatalf("Token nicht gesendet")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// token = strings.Replace(token, "Bearer ", "", 1)

	UserInfo, err := Get_UserID(token)

	if err != nil {
		log.Default().Printf(err.Error())
		return c.SendStatus(401)
	}

	// Oauth login -> new Token
	newToken, err := Get_AdminToken()

	if err != nil {
		log.Default().Printf(err.Error())
		return c.SendStatus(400)
	}

	userGroupes, err := Get_UserGroups(newToken, UserInfo.Sub)

	if err != nil {
		log.Default().Printf(err.Error())
		return c.SendStatus(416)
	}

	userGroups := []string{
		"Takelmeister",
		"Admin",
	}

	if Check_IsUserPartOfGroup(userGroups, userGroupes) {
		return c.Next()
	}

	return c.SendStatus(401)
}

func Check_IsUserLoggedIn(c *fiber.Ctx) error {

	admingroups := []string{
		"Takelmeister",
		"Entwickler",
		"Admin",
	}

	// First Check for arbeitsstundenToken (Token of this service)
	// ArbeitsstundenCooki
	firstTokenValue := c.Cookies("ArbeitsstundenCooki")

	if firstTokenValue != "" {
		vallid, info := CheckCookie(firstTokenValue)

		if vallid {
			if Check_IsUserPartOfGroup(admingroups, info.Groups) {
				c.Append("isTakel", "true")
				c.Cookie(&fiber.Cookie{
					Name:  "isTakel",
					Value: "true",
				})
			}
			return c.Next()
		}
	}

	// else, check for keycloak Token
	token := c.Cookies("token")
	if token == "" {
		log.Fatalf("kein Token gesendet")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// token = strings.Replace(token, "Bearer ", "", 1)

	var currentMember gocloak.User
	currentUser, err := Get_UserID(token)
	groups, groupErr := Get_UserGroups(token, currentUser.Sub)

	if err != nil {
		log.Default().Printf(err.Error())
		return c.SendStatus(401)
	}
	if groupErr != nil {
		log.Default().Printf(err.Error())
		return c.SendStatus(401)
	}

	currentMember.Info = currentUser
	currentMember.Groups = groups

	if Check_IsUserPartOfGroup(admingroups, groups) {
		c.Append("isTakel", "true")
		c.Cookie(&fiber.Cookie{
			Name:  "isTakel",
			Value: "true",
		})
	}

	newCookie := CreateCookie(currentMember)
	c.Cookie(&newCookie)

	return c.Next()
}
