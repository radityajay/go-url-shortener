package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/pandeptwidyaop/golog"
	"github.com/radityajay/go-url-shortener/db"
	"github.com/radityajay/go-url-shortener/routes"
	"github.com/radityajay/go-url-shortener/utils/logc"
	"github.com/radityajay/go-url-shortener/utils/wg"
)

func main() {
	appName := os.Getenv("APP_NAME")
	listenPort := ":4000"

	InitEnv()
	wg.NewHttpWg()
	InitDatabase()
	golog.New()

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(
		logger.New(logger.Config{
			Format:     "${time} | ${green} ${status} ${white} | ${latency} | ${ip} | ${green} ${method} ${white} | ${path} | ${yellow} ${body} ${reset} | ${magenta} ${resBody} ${reset}\n",
			TimeFormat: "02 January 2006 15:04:05",
			TimeZone:   "Asia/Jakarta",
		}),

		func(c *fiber.Ctx) error {
			defer func() {
				if r := recover(); r != nil {
					dbg := debug.Stack()
					logc.Error("Server panic Occured", fmt.Errorf("%s", r), &dbg)
					c.Status(fiber.StatusInternalServerError).SendString("Server Error")
				}
			}()
			return c.Next()
		},
	)

	// Init Routing
	routes.RouteApiRegister(app)

	go func() {
		golog.Slack.Info(fmt.Sprintf("%s: HTTP-Only Service Started", appName))
		log.Fatal(app.Listen(listenPort))
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	wg.HttpWG.Wait()

	golog.Slack.Info(fmt.Sprintf("%s: HTTP-Only Service Stopped", appName))
	app.Shutdown()
}

func InitEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found")
	}
}

func InitDatabase() {
	db.NewPostgres()
}
