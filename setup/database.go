package setup

import (
	"os"

	db "github.com/ASV-Aachen/Seereisenplan-backend/modules/DB"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MariaDB
var DB_USER string = os.Getenv("DB_USER")
var DB_PASSWORD string = os.Getenv("DB_PASSWORD")
var DB_NAME string = os.Getenv("DB_NAME")
var DB_URL string = os.Getenv("DB_URL")

func DB_Migrate(selectedDatabase *gorm.DB) {
	// Migrate the schema
	err := selectedDatabase.AutoMigrate(
		&db.User{},
		&db.Season{},
		&db.Reduction{},
		&db.Project{},
		&db.Project_item{},
		&db.Project_item_hour{},
	)

	if err != nil {
		panic(err.Error())
	}
}

func SetUpMariaDB() *gorm.DB {
	dsn := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_URL + ":3306" + ")/" + DB_NAME

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db
}
