package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	// CSRF setup
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "form:csrf",
		CookieName:     "csrf_",
		ContextKey:     "csrf", // key to retrieve generated token
		CookieSameSite: "Lax",
	}))

	app = setupApp()

	log.Fatal(app.Listen(":3000"))
}

func setupApp () *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	// CSRF setup
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "form:csrf",
		CookieName:     "csrf_",
		ContextKey:     "csrf", // key to retrieve generated token
		CookieSameSite: "Lax",
	}))


	// GET form
	app.Get("/form", func(c *fiber.Ctx) error {
		// Retrieve token added by middleware
		token := c.Locals("csrf")
		if token == nil {
			return fiber.NewError(fiber.StatusInternalServerError, "CSRF token missing")
		}

		html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Fiber CSRF</title></head>
<body>
    <h1>Fiber CSRF Form</h1>
    <form action="/submit" method="POST">
        <label>Name: <input type="text" name="name"></label>
        <input type="hidden" name="csrf" value="%s">
        <button type="submit">Submit</button>
    </form>
</body>
</html>
`, token)

		return c.Type("html").SendString(html)
	})

	// POST form
	app.Post("/submit", func(c *fiber.Ctx) error {
		// If CSRF token is bad, Fiber auto-403s before this handler runs
		name := c.FormValue("name", "anonymous")
		return c.SendString(fmt.Sprintf("Hello %s! CSRF validation succeeded.\n", name))
	})

	return app

}