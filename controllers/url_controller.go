package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radityajay/go-url-shortener/utils/logc"
	"github.com/radityajay/go-url-shortener/utils/url"
	"github.com/radityajay/go-url-shortener/utils/wg"
)

func CreateShortURL(c *fiber.Ctx) error {
	wg.HttpWG.Add(1)
	defer wg.HttpWG.Done()

	type RequestShortenerURL struct {
		LongURL     string  `json:"long_url"`
		CustomAlias *string `json:"custom_alias"`
	}

	req := new(RequestShortenerURL)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	errors := ValidateRequest(req)

	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  false,
			"message": "Data not valid",
			"error":   errors,
		})
	}

	res, err := url.CreateGenerateURL(req.LongURL, req.CustomAlias)

	if err != nil {
		var dataReport []byte = c.Body()
		switch err.Error() {
		case "failed to create":
			logc.ErrorWithFormat("creating url", logc.ERR, err, &dataReport)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "error when creating url",
			})
		default:
			logc.ErrorWithFormat("creating url", logc.UNHANDLED, err, &dataReport)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "unknown error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"url":    res,
		"message": "Success",
	})
}

func GetShortURL(c *fiber.Ctx) error {
	wg.HttpWG.Add(1)
	defer wg.HttpWG.Done()

	url, err := url.GetRedirectURL(c.Params("short_code"))

	if err != nil {
		var dataReport []byte = c.Body()
		switch err.Error() {
		case "not found":
			logc.ErrorWithFormat("redirect url", logc.ERR, err, &dataReport)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "url not found",
			})
		default:
			logc.ErrorWithFormat("redirect url", logc.UNHANDLED, err, &dataReport)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "unknown error",
			})
		}
	}

	return c.Redirect(url, fiber.StatusMovedPermanently)
}
