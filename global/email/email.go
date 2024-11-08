package email

import (
	"errors"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	return godotenv.Load()
}

func SendEmailSMTP(to string, emailBody string, verificationCode string) (bool, error) {
	emailID := os.Getenv("EMAIL_ID")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	if emailID == "" || emailPassword == "" {
		return false, errors.New("email credentials are not set in environment variables")
	}

	emailAuth := smtp.PlainAuth("", emailID, emailPassword, "smtp.gmail.com")

	subject := "Subject: GOMS 인증번호\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + mime + "\n" + emailBody)

	if err := smtp.SendMail("smtp.gmail.com:587", emailAuth, emailID, []string{to}, msg); err != nil {
		return false, err
	}

	return true, nil
}
