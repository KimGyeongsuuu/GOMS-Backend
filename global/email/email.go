package email

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	return godotenv.Load()
}

func SendEmailSMTP(to string, data interface{}, templateName string, verificationCode string) (bool, error) {
	emailID := os.Getenv("EMAIL_ID")
	emailPassword := os.Getenv("EMAIL_PASSWORD")

	if emailID == "" || emailPassword == "" {
		return false, errors.New("email credentials are not set in environment variables")
	}

	emailAuth := smtp.PlainAuth("", emailID, emailPassword, "smtp.gmail.com")

	emailData := map[string]interface{}{
		"VerificationCode": verificationCode,
	}

	emailBody, err := parseTemplate(templateName, emailData)
	if err != nil {
		return false, err
	}

	subject := "Subject: GOMS 인증번호\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + mime + "\n" + emailBody)

	if err := smtp.SendMail("smtp.gmail.com:587", emailAuth, emailID, []string{to}, msg); err != nil {
		return false, err
	}

	return true, nil
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath := fmt.Sprintf("email_templates/%s", templateFileName)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
