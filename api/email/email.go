// api/email/email.go
package email

import (
	"fmt"
	"net/smtp"
)

// SmtpConfig estructura para configuración SMTP
type SmtpConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// SendMail función para enviar correo electrónico usando SMTP
func SendMail(smtpConfig SmtpConfig, to, subject, body string) error {
	fmt.Println("Sending email to:", to)
	// Flag para depurar el servidor SMTP (imprimir en consola) LOS DATOS DE ENVÍO
	fmt.Println("Host:", smtpConfig.Host, "Port:", smtpConfig.Port, "Username:", smtpConfig.Username, "Password:", smtpConfig)
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpConfig.Host+":"+smtpConfig.Port, auth, smtpConfig.Username, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}
