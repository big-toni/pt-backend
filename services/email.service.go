package services

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailService struct
type EmailService struct {
	client *sendgrid.Client
}

// NewEmailService creates a new EmailService.
func NewEmailService() *EmailService {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	client := sendgrid.NewSendClient(apiKey)
	return &EmailService{client}
}

// SendResetEmail func
func (s *EmailService) SendResetEmail(email string, url string) (bool, error) {
	from := mail.NewEmail("PT Support", "test@example.com")
	subject := "Reset Password"
	to := mail.NewEmail(email, email)
	plainTextContent := "Press the link to reset password: " + url
	htmlContent := `<strong>Press the link to reset password: </strong>
	<a href=` + url + `>Click here</a> mean?`

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	res, err := s.client.Send(message)

	if err != nil {
		log.Println(err)
		return false, err
	}

	fmt.Println(res.StatusCode)
	fmt.Println(res.Body)
	fmt.Println(res.Headers)

	return true, nil

}
