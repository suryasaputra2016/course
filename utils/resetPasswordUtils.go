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

	auth := smtp.PlainAuth("", username, password, host)
	subject := "Reset Password"
	from := mail.Address{
		Name:    "admin",
		Address: "admin@course.com",
	}
	to := mail.Address{
		Name:    "mr./mrs.",
		Address: email,
	}
	htmlBody := "<h1>Reset password link</h1><p>Link: <a href=\"#\">" + token + "</a></p>"

	headers := map[string]string{
		"Subject":      subject,
		"From":         from.String(),
		"To":           to.String(),
		"MIME-version": "1.0;",
		"Content-Type": "text/html; charset=\"UTF-8\";",
	}

	var message string
	for k, v := range headers {
		message += k + ": " + v + "\n"
	}
	message += "\n" + htmlBody

	err := smtp.SendMail(address, auth, from.Address, []string{to.Address}, []byte(message))
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
