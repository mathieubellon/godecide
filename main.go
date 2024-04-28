package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"github.com/shareed2k/goth_fiber"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	connectDB()
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	app.Use(logger.New())
	app.Use(csrf.New())
	app.Use(recover.New())
	app.Use(helmet.New())
	// 404 Handler

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:8000, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Static("/static", "./static")

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://127.0.0.1:8000/auth/callback/google"), // TODO make BASE_URL an env variable
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://127.0.0.1:8000/auth/callback/github"),
	)

	app.Get("/", Homepage).Name("index") // Serves vue frontend
	app.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	app.Get("/auth/callback/:provider", Callback)
	app.Get("/logout", Logout)

	api := app.Group("/api") // /api
	api.Use(Protected)
	api.Get("/me", Me)
	v1 := api.Group("/v1") // /api/v1
	v1.Get("/ideas", ListIdeas)

	// Catch All route for 404 : TODO : handle this frontend only ?
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
	// data, _ := json.MarshalIndent(app.GetRoutes(true), "", "  ")
	// fmt.Print(string(data))

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}
