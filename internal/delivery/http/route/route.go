package route

import (
	"todo-app/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
	RoleController *http.RoleController
	AuthMiddleware fiber.Handler
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

	c.App.Get("/api/roles", c.RoleController.List)
	c.App.Post("/api/roles", c.RoleController.Create)
	c.App.Put("/api/roles/update/:roleId", c.RoleController.Update)
	c.App.Get("/api/roles/view/:roleId", c.RoleController.Get)
	c.App.Put("/api/roles/delete/:roleId", c.RoleController.SoftDelete)
	c.App.Get("/api/roles/trash", c.RoleController.RecycleBin)
	c.App.Put("/api/roles/restore/:roleId", c.RoleController.Restore)
	c.App.Delete("/api/roles/force/:roleId", c.RoleController.ForceDelete)

}
