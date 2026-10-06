package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(to string, token string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	appURL := os.Getenv("APP_URL")

	verifyLink := fmt.Sprintf("%s/verify-email?token=%s", appURL, token)

	subject := "Subject: Verify your email\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	body := fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Please click the link below to verify your email:</p>
		<a href="%s">Verify Email</a>
		<br><br>
		<p>Or copy this link:</p>
		<p>%s</p>
	`, verifyLink, verifyLink)

	message := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		smtpEmail,
		[]string{to},
		message,
	)
}



func SendPasswordResetEmail(to string, token string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	frontendURL := os.Getenv("APP_URL")

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
    
    fmt.Println(frontendURL)
    fmt.Println(resetLink)

	subject := "Subject: Reset your password\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	body := fmt.Sprintf(`
		<h2>Password Reset</h2>

		<p>Click the link below to reset your password:</p>
		<p><a href="%s">Reset Password</a></p>

		<p><strong>Reset Token:</strong></p>
		<p style="word-break:break-all;">%s</p>

		<p><strong>Direct API JSON:</strong></p>
		<pre>
{
  "token":"%s",
  "new_password":"YourNewPassword123"
}
		</pre>

		<p>This link will expire in 15 minutes.</p>
		<p>If you did not request this, ignore this email.</p>
	`, resetLink, token, token)

	message := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		smtpEmail,
		[]string{to},
		message,
	)
}