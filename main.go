package main

import (
	"github.com/ASV-Aachen/Seereisenplan-backend/cmd"
	"github.com/ASV-Aachen/Seereisenplan-backend/setup"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	// Test handler
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("App running")
	})

	var db *gorm.DB = setup.SetUpMariaDB()
	setup.DB_Migrate(db)

	api := app.Group("/seereisenplan/V1/api")

	api.Use(
		setup.Check_IsUserLoggedIn,
	)

	app.Get("/", func(c *fiber.Ctx) error { return cmd.GetCruisesForCurrentYear(c, db) })
	app.Get("/:year", func(c *fiber.Ctx) error { return cmd.GetCruisesForYear(c, db) })

	app.Get("Licenses", func(c *fiber.Ctx) error { return cmd.GetLicenses(c, db) })

	users := app.Group("/sailor")
	app.Get("/", func(c *fiber.Ctx) error { return cmd.GetAllUsers(c, db) })

	user := users.Group("/:userID")
	user.Get("/", func(c *fiber.Ctx) error { return cmd.GetUser(c, db) })
	user.Post("/", func(c *fiber.Ctx) error { return cmd.NewUser(c, db) })
	user.Patch("/", func(c *fiber.Ctx) error { return cmd.UpdateUser(c, db) })

	app.Post("/", func(c *fiber.Ctx) error { return cmd.NewCruise(c, db) })
	cruise := app.Group("/cruise/:cruiseID")
	cruise.Get("/", func(c *fiber.Ctx) error { return cmd.GetCruise(c, db) })
	cruise.Delete("/", func(c *fiber.Ctx) error { return cmd.RemoveCruise(c, db) })
	cruise.Patch("/", func(c *fiber.Ctx) error { return cmd.UpdateCruise(c, db) })

	app.Listen(":3000")
}
