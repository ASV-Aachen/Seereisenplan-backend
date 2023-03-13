package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ASV-Aachen/Seereisenplan-backend/modules/db"
	"github.com/ASV-Aachen/Seereisenplan-backend/modules/json_struct"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetCruisesForCurrentYear(c *fiber.Ctx, database *gorm.DB) error {
	var AllCruises []db.Cruise

	currentYear := time.Now().Year()
	database.Model(&db.Cruise{Season: currentYear}).Find(AllCruises)

	var result []json_struct.Cruises

	for _, s := range AllCruises {
		temp := new(json_struct.Cruises)
		temp.CruiseName = s.CruiseName
		temp.CuriseDescription = s.CuriseDescription
		temp.StartDate = s.StartDate
		temp.EndDate = s.EndDate
		temp.StartPort = s.StartPort
		temp.EndPort = s.EndPort

		var Share []db.CruiseShare
		database.Model(&db.CruiseShare{Cruise: s}).Find(Share)

		for _, t := range Share {
			tempsailor := new(json_struct.Sailor)
			tempsailor.First_name = t.Sailor.First_name
			tempsailor.Last_name = t.Sailor.Last_name
			tempsailor.Position = t.Position

			temp.Sailor = append(temp.Sailor, *tempsailor)
		}

	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func GetCruisesForYear(c *fiber.Ctx, database *gorm.DB) error {
	var AllCruises []db.Cruise

	currentYear, err := strconv.Atoi(c.Params("year"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	database.Model(&db.Cruise{Season: currentYear}).Find(AllCruises)

	var result []json_struct.Cruises

	for _, s := range AllCruises {

		temp := new(json_struct.Cruises)
		temp.CruiseName = s.CruiseName
		temp.CuriseDescription = s.CuriseDescription
		temp.StartDate = s.StartDate
		temp.EndDate = s.EndDate
		temp.StartPort = s.StartPort
		temp.EndPort = s.EndPort

		var Share []db.CruiseShare
		database.Model(&db.CruiseShare{Cruise: s}).Find(Share)

		for _, t := range Share {
			tempsailor := new(json_struct.Sailor)
			tempsailor.First_name = t.Sailor.First_name
			tempsailor.Last_name = t.Sailor.Last_name
			tempsailor.Position = t.Position

			temp.Sailor = append(temp.Sailor, *tempsailor)
		}

	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func GetCruise(c *fiber.Ctx, database *gorm.DB) error {
	searched_id := c.Params("cruiseID")
	var s db.Cruise

	database.Model(&db.Cruise{ID: searched_id}).First(s)

	var result json_struct.Cruises

	temp := new(json_struct.Cruises)
	temp.CruiseName = s.CruiseName
	temp.CuriseDescription = s.CuriseDescription
	temp.StartDate = s.StartDate
	temp.EndDate = s.EndDate
	temp.StartPort = s.StartPort
	temp.EndPort = s.EndPort

	var Share []db.CruiseShare
	database.Model(&db.CruiseShare{Cruise: s}).Find(Share)

	for _, t := range Share {
		tempsailor := new(json_struct.Sailor)
		tempsailor.First_name = t.Sailor.First_name
		tempsailor.Last_name = t.Sailor.Last_name
		tempsailor.Position = t.Position

		temp.Sailor = append(temp.Sailor, *tempsailor)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func NewCruise(c *fiber.Ctx, database *gorm.DB) error {
	errorAnswer := ""
	json := new(json_struct.JSONCruise)
	if err := c.BodyParser(json); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	tempConnectStatus := database.Create(&db.Cruise{
		CruiseName:        json.CruiseName,
		CuriseDescription: json.CuriseDescription,
		StartDate:         json.StartDate,
		EndDate:           json.EndDate,
		StartPort:         json.StartPort,
		EndPort:           json.EndPort,
	})

	if tempConnectStatus.Error != nil {
		errorAnswer += fmt.Sprintf("Reise %s konnte nicht angelegt werden", json.CruiseName)
		errorAnswer += tempConnectStatus.Error.Error() + "\n"
	}

	var currentCruise db.Cruise
	tempConnectStatus.First(currentCruise)
	// database.Model(&db.Cruise{ID: }).First(currentCruise)

	for _, i := range json.JSONSailor {
		errorAnswer += "\n" + addCruiseShare(database, i, currentCruise, errorAnswer)
	}

	if errorAnswer != "" {
		return c.Status(fiber.StatusConflict).SendString(errorAnswer)
	}

	return c.Status(fiber.StatusOK).SendString("Angelegt")
}

func addCruiseShare(database *gorm.DB, i json_struct.JSONSailor, currentCruise db.Cruise, errorAnswer string) string {
	var currentSailor db.Sailor
	database.Model((&db.Sailor{ID: i.ID})).First(currentSailor)

	tempSailorStatus := database.Create(&db.CruiseShare{
		Sailor:   currentSailor,
		Cruise:   currentCruise,
		Position: i.Position,
	})
	if tempSailorStatus.Error != nil {
		errorAnswer += fmt.Sprintf("Segler %s konnte zur Reise %s konnte nicht angelegt werden", currentSailor.Last_name)
		errorAnswer += tempSailorStatus.Error.Error() + "\n"
	}
	return errorAnswer
}

func RemoveCruise(c *fiber.Ctx, database *gorm.DB) error {
	searched_id := c.Params("cruiseID")

	database.Delete(&db.Cruise{}, searched_id)

	return c.Status(fiber.StatusOK).SendString("removed")
}

func UpdateCruise(c *fiber.Ctx, database *gorm.DB) error {
	searched_id := c.Params("cruiseID")
	errorAnswer := ""
	json := new(json_struct.JSONCruise)
	if err := c.BodyParser(json); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var currentCruise db.Cruise
	database.Model(&db.Cruise{ID: searched_id}).First(currentCruise)

	currentCruise.CruiseName = json.CruiseName
	currentCruise.CuriseDescription = json.CuriseDescription
	currentCruise.StartDate = json.StartDate
	currentCruise.EndDate = json.EndDate
	currentCruise.StartPort = json.StartPort
	currentCruise.EndPort = json.EndPort

	var CurrentShares []db.CruiseShare
	database.Model(&db.CruiseShare{Cruise: currentCruise}).Find(CurrentShares)

	// Check für neue Mitglieder und leg an wenn nicht da
	for _, jsonsailor := range json.JSONSailor {
		for _, share := range CurrentShares {
			if share.Sailor.ID == jsonsailor.ID {
				continue
			}
			errorAnswer += addCruiseShare(database, jsonsailor, currentCruise, errorAnswer)
		}
	}

	// Check für alte Mitglieder und lösch die Raus
	for _, share := range CurrentShares {
		for _, jsonsailor := range json.JSONSailor {
			if share.Sailor.ID == jsonsailor.ID {
				continue
			}
		}
		database.Delete(&db.CruiseShare{}, share)
	}

	database.Save(currentCruise)

	return c.Status(fiber.StatusOK).SendString("OK")
}
