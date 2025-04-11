package utils

import (
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func CheckEmailFormat(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re, err := regexp.Compile(emailRegex)
	if err != nil {

	}
	if re.MatchString(email) {
		return nil
	} else {
		return errors.New("email is not well formatted")
	}
}

func SendPasswordResetEmail(email, token string) error {
	godotenv.Load()
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	address := os.Getenv("ADDRESS")

	from := mail.Address{
		Name:    "admin",
		Address: "admin@course.com",
	}

	to := mail.Address{
		Name:    "name",
		Address: email,
	}

	auth := smtp.PlainAuth("", username, password, host)
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\nFrom: " + from.String() + "\nTo: " + to.String()
	subject := "Reset Password"
	htmlBody := "<h1>Reset password link</h1><p>Link: <a href=\"#\">click here: " + token + "</a></p>"
	message := "Subject: " + subject + "\n" + headers + "\n\n" + htmlBody

	err := smtp.SendMail(address, auth, from.Address, []string{to.Address}, []byte(message))
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
