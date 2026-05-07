package helper

import (
	"fmt"
	"net/smtp"

	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/config"
)

func SendEmail(cfg config.SMTPConfig, to string, subject string, body string) error {
	header := make(map[string]string)
	header["From"] = cfg.Email
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", cfg.Email, cfg.Password, cfg.Host)

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return smtp.SendMail(address, auth, cfg.Email, []string{to}, []byte(message))
}
