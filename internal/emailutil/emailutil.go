package emailutil

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/dagota12/Loan-Tracker/bootstrap"
)

// SendVerificationEmail sends an account verification email to the user.
func SendVerificationEmail(recipientEmail string, VerificationToken string, env *bootstrap.Env) error {
	// Email configuration
	from := env.SenderEmail
	password := env.SenderPassword
	smtpHost := env.SmtpHost
	smtpPort := env.SmtpPort

	subject := "Subject: Account Verification\n"
	mime := "MIME-Version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	url := fmt.Sprintf("http://localhost:8080/users/verify-email/%v", VerificationToken)
	body := EmailTemplate(url) // Assuming Emailtemplate(url) returns a string containing the HTML email body
	message := []byte(subject + mime + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	log.Printf("Sending email to %s using SMTP server %s:%s", recipientEmail, smtpHost, smtpPort)

	// Correcting the SendMail function to use the actual message content
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipientEmail}, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Verification email sent successfully.")
	return nil
}

func SendOtpVerificationEmail(recipientEmail string, otp string, env *bootstrap.Env) error {
	// Email configuration
	from := env.SenderEmail
	password := env.SenderPassword
	smtpHost := env.SmtpHost
	smtpPort := env.SmtpPort

	subject := "Subject: Account Verification\n"
	mime := "MIME-Version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := OTPEmailTemplate(otp, env)
	message := []byte(subject + mime + "\n" + body)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipientEmail}, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	return nil

}
