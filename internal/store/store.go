package store

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StoreLayer interface {
	CreateSubmission(formId uint, ipAddress string, userAgent string, values map[uint]string) (Submission, error)
	CreateSubmissionValue(submissionId uint, formValueId uint, value string) (SubmissionValue, error)
	GetForm(slug string) (Form, error)
	GetFormExists(slug string) bool
	GetFormField(formID uint, slug string) (FormField, error)
	GetFormFields(formID uint) ([]FormField, error)
}

type storeLayer struct {
	db *gorm.DB
}

func New() *storeLayer {
	username, usernameExists := os.LookupEnv("DB_USERNAME")
	if !usernameExists {
		log.Fatal("Error: No database username provided")
		return nil
	}

	password, passwordExists := os.LookupEnv("DB_PASSWORD")
	if !passwordExists {
		log.Fatal("Error: No database password provided")
		return nil
	}

	host, hostExists := os.LookupEnv("DB_HOST")
	if !hostExists {
		log.Fatal("Error: No database hostname provided")
		return nil
	}

	name, nameExists := os.LookupEnv("DB_NAME")
	if !nameExists {
		log.Fatal("Error: No database name provided")
		return nil
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, name)), &gorm.Config{})
	if err != nil {
		log.Fatal("Error: Unable to connect to MySQL server", err)
	}

	db.AutoMigrate(&Form{})
	db.AutoMigrate(&FormField{})
	db.AutoMigrate(&Submission{})
	db.AutoMigrate(&SubmissionValue{})

	return &storeLayer{
		db: db,
	}
}
