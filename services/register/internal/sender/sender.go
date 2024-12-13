package sender

import (
	"fmt"
	"log/slog"
	"net/smtp"
)

type Config struct {
	Email    string `yaml:"email" env-prefix:"EMAIL"`
	Password string `yaml:"password" env-prefix:"PASSWORD"`
	Host     string `yaml:"host" env-prefix:"HOST"`
	Port     string `yaml:"port" env-prefix:"PORT"`
}

type EmailSender struct {
	cfg  *Config
	auth *smtp.Auth
}

func New(cfg *Config) (*EmailSender, error) {
	auth := smtp.PlainAuth("", cfg.Email, cfg.Password, cfg.Host)
	sender := EmailSender{cfg: cfg, auth: &auth}
	return &sender, nil
}

func (s *EmailSender) SendEmail(content []byte, reciever string) error {
	slog.Info(fmt.Sprintf("Sending code %s to user %s", string(content), reciever))
	err := smtp.SendMail(s.cfg.Host+":"+s.cfg.Port, *s.auth, s.cfg.Email, []string{reciever}, content)
	if err != nil {
		slog.Error(fmt.Sprintf("sender SendEmail error: %v", err))
	}
	return err
}
