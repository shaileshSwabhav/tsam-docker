package service

import (
	"net/smtp"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/config"
	contactCfg "github.com/techlabs/swabhav/tsam/homePage/contact"
	"github.com/techlabs/swabhav/tsam/models/homePage/contact"
	"github.com/techlabs/swabhav/tsam/repository"
)

type ContactInfoService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

func NewContactInfoService(db *gorm.DB, repo repository.Repository) *ContactInfoService {
	return &ContactInfoService{
		DB:         db,
		Repository: repo,
	}
}

func (service *ContactInfoService) ContactUs(content *contact.ContactInfo, email, pass string) error {
	msgs, err := content.CreateMessageBody()
	if err != nil {
		return err
	}
	// return nil
	from := email
	password := pass

	// smtp server configuration.
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	// Message.
	message := []byte(msgs)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, contactCfg.RecipientMails, message)
	if err != nil {
		return err
	}
	return nil
}
