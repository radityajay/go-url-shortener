package controllers

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type ValidationResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func ValidateRequest(toValidate interface{}) *map[string]string {
	var msgs = map[string]string{}
	success, err := govalidator.ValidateStruct(toValidate)
	if !success {
		ctx := strings.Split(err.Error(), ";")
		for _, e := range ctx {
			msg := strings.Split(e, ": ")
			msgs[msg[0]] = msg[1]
		}
	}
	if !success {
		return &msgs
	} else {
		return nil
	}
}

func Default(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
	})
}
