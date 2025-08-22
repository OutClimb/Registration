package store

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StoreLayer interface {
	CreateForm(name, slug, template string, opensOn, closesOn *time.Time, maxSubmissions uint, notOpenMessage, closedMessage, successMessage *string, emailFormFieldSlug, emailTo, emailSubject, emailTemplate string) (*Form, error)
	CreateFormField(formID uint, name, slug, fieldType string, metadata *string, required bool, validation *string, order uint) (*FormField, error)
	CreateSubmission(form *Form, fields *[]FormField, ipAddress string, userAgent string, values map[string]string) (*Submission, error)
	CreateSubmissionValue(submissionId uint, formValueId uint, value string) (*SubmissionValue, error)
	DeleteSubmission(id uint) error
	GetAllForms() (*[]Form, error)
	GetForm(slug string) (*Form, error)
	GetFormFields(formID uint) (*[]FormField, error)
	GetFormsForUser(userId uint) (*[]Form, error)
	GetNumberOfSubmissions(formID uint) (int64, error)
	GetSubmissionsForForm(formID uint) (*[]Submission, error)
	GetSubmissionValues(submissionId uint) (*[]SubmissionValue, error)
	GetUser(id uint) (*User, error)
	GetUserWithUsername(username string) (*User, error)
	UpdatePassword(id uint, password string) error
}

type storeLayer struct {
	db *gorm.DB
}

func New() *storeLayer {
	username := os.Getenv("DB_USERNAME")
	if len(username) == 0 {
		log.Fatal("Error: No database username provided")
		return nil
	}

	password := os.Getenv("DB_PASSWORD")

	host := os.Getenv("DB_HOST")
	if len(host) == 0 {
		log.Fatal("Error: No database hostname provided")
		return nil
	}

	name := os.Getenv("DB_NAME")
	if len(name) == 0 {
		log.Fatal("Error: No database name provided")
		return nil
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, name)), &gorm.Config{})
	if err != nil {
		log.Fatal("Error: Unable to connect to MySQL server", err)
	}

	db.AutoMigrate(&Form{})
	db.AutoMigrate(&FormField{})
	db.AutoMigrate(&Submission{})
	db.AutoMigrate(&SubmissionValue{})
	db.AutoMigrate(&User{})

	return &storeLayer{
		db: db,
	}
}
