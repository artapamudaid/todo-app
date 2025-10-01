package route

import (
	"todo-app/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	UserController    *http.UserController
	ContactController *http.ContactController
	AddressController *http.AddressController
	AuthMiddleware    fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/auth/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Delete("/api/auth/logout", c.UserController.Logout)
	c.App.Patch("/api/profile/update", c.UserController.Update)
	c.App.Get("/api/profile", c.UserController.Current)
	c.App.Post("/api/auth/refresh", c.UserController.Refresh)

	c.App.Get("/api/contacts", c.ContactController.List)
	c.App.Post("/api/contacts", c.ContactController.Create)
	c.App.Put("/api/contacts/:contactId", c.ContactController.Update)
	c.App.Get("/api/contacts/:contactId", c.ContactController.Get)
	c.App.Delete("/api/contacts/:contactId", c.ContactController.Delete)

	c.App.Get("/api/contacts/:contactId/addresses", c.AddressController.List)
	c.App.Post("/api/contacts/:contactId/addresses", c.AddressController.Create)
	c.App.Put("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Update)
	c.App.Get("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Get)
	c.App.Delete("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Delete)
}
