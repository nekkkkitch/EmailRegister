package router

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Router struct {
	cfg     *Config
	app     *fiber.App
	service IService
}

type IService interface {
	Register(email, password string) error
	VerifyEmail(email, code string) error
}

func New(cfg *Config, service IService) (*Router, error) {
	app := fiber.New()
	router := Router{cfg: cfg, app: app, service: service}
	router.app.Post("/register", router.Register())
	router.app.Post("/verifyemail", router.VerifyEmail())
	return &router, nil
}

func (r *Router) Listen() error {
	addr := r.cfg.Host + ":" + r.cfg.Port
	err := r.app.Listen(addr)
	if err != nil {
		log.Printf("Failed to listen on %v: %v", addr, err)
		return err
	}
	return nil
}

func (r *Router) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := map[string]string{}
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Println("Failed to unmarshal register body:", err)
			c.Status(500)
			return nil
		}
		err = r.service.Register(data["email"], data["password"])
		if err != nil {
			log.Println("Failed to register user:", err)
			c.Status(500)
			return nil
		}
		return nil
	}
}

func (r *Router) VerifyEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := map[string]string{}
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Println("Failed to unmarshal verify body:", err)
			c.Status(500)
			return nil
		}
		err = r.service.VerifyEmail(data["email"], data["code"])
		if err != nil {
			if err.Error() == "codes not equal" {
				c.Status(fiber.StatusBadRequest)
				return err
			}
			log.Println("Failed to verify user:", err)
			c.Status(500)
			return nil
		}
		return nil
	}
}
