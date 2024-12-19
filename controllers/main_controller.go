package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radityajay/go-url-shortener/utils/wg"
)

func Main(c *fiber.Ctx) error {
	wg.HttpWG.Add(1)
	defer wg.HttpWG.Done()

	return c.SendString("OK")
}
