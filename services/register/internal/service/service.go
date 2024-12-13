package service

import (
	"emailregister/pkg/crypt"
	"fmt"
	"log"
	"math/rand/v2"
)

type Service struct {
	redis  IRedis
	db     IDB
	sender ISender
}

type IRedis interface {
	GetCode(email string) (string, error)
	PutCode(email, code string) error
	DeleteCode(email string) error
}

type IDB interface {
	AddUser(email, password string) error
	SetUserVerificationStatus(email string, status bool) error
}

type ISender interface {
	SendEmail(content []byte, reciever string) error
}

func New(redis IRedis, db IDB, sender ISender) *Service {
	service := Service{redis: redis, db: db, sender: sender}
	return &service
}

func (s *Service) Register(email, password string) error {
	hashedPass, err := crypt.CryptPassword(password)
	if err != nil {
		log.Println("Failed crypting password:", err)
		return err
	}
	err = s.db.AddUser(email, string(hashedPass))
	if err != nil {
		log.Println("Failed adding user to db:", err)
		return err
	}
	verificationCode := generateCode()
	err = s.redis.PutCode(email, verificationCode)
	if err != nil {
		log.Println("Failed saving verification code:", err)
		return err
	}
	message := "Здравствуйте! Ваш код верификации: " + verificationCode
	err = s.sender.SendEmail([]byte(message), email)
	if err != nil {
		log.Println("Failed sending verification code:", err)
		return err
	}
	return nil
}

func (s *Service) VerifyEmail(email, code string) error {
	sentCode, err := s.redis.GetCode(email)
	if err != nil {
		log.Println("Failed getting code from redis:", err)
		return err
	}
	if sentCode != code {
		log.Println("User entered wrong code")
		return fmt.Errorf("codes not equal")
	}
	err = s.db.SetUserVerificationStatus(email, true)
	if err != nil {
		log.Println("Failed to set user verification status:", err)
		return err
	}
	err = s.redis.DeleteCode(email)
	if err != nil {
		log.Println("Failed deleting code:", err)
		return nil
	}
	return nil
}

func generateCode() string {
	n := 6
	code := ""
	alph := "1234567890"
	for range n {
		code += string([]rune(alph)[rand.IntN(len(alph))])
	}
	return code
}
