package cmd

import (
	"github.com/ASV-Aachen/Seereisenplan-backend/modules/db"
	"github.com/ASV-Aachen/Seereisenplan-backend/modules/json_struct"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllUsers(c *fiber.Ctx, database *gorm.DB) error {
	var currentSailors []db.Sailor
	database.Model(&db.Sailor{}).Find(currentSailors)

	return c.Status(fiber.StatusOK).JSON(currentSailors)
}

func GetUser(c *fiber.Ctx, database *gorm.DB) error {
	searched_id := c.Params("userID")
	var currentSailors db.Sailor
	database.Model(&db.Sailor{ID: searched_id}).First(currentSailors)

	return c.Status(fiber.StatusOK).JSON(currentSailors)
}

func NewUser(c *fiber.Ctx, database *gorm.DB) error {
	// Anlegen von Gast und nur von Gast
	json := new(json_struct.JSON_NEW_GUEST)
	if err := c.BodyParser(json); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	database.Create(&db.Sailor{
		First_name: json.First_name,
		Last_name:  json.Last_name,
		Guest:      true,
	})

	return c.Status(fiber.StatusOK).SendStatus(200)
}

func remove(s []db.LicenseShare, i int) []db.LicenseShare {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func UpdateUser(c *fiber.Ctx, database *gorm.DB) error {
	searched_id := c.Params("userID")

	var json []db.LicenseShare
	if err := c.BodyParser(json); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var currentSailors db.Sailor
	database.Model(&db.Sailor{ID: searched_id}).First(currentSailors)

	for _, i := range json {
		for _, currentLicense := range currentSailors.License {
			if i == currentLicense {
				continue
			}
		}
		// anlegen
		currentSailors.License = append(currentSailors.License, i)
	}

	for x, i := range currentSailors.License {
		for _, jsonLicesne := range json {
			if i == jsonLicesne {
				continue
			}
		}
		// löschen
		remove(currentSailors.License, x)
	}

	database.Save(currentSailors)
	return c.Status(fiber.StatusOK).SendStatus(200)
}

func GetLicenses(c *fiber.Ctx, database *gorm.DB) error {
	var AllLicenses []db.License
	database.Model(&db.License{}).Find(AllLicenses)

	return c.Status(fiber.StatusOK).JSON(AllLicenses)
}
