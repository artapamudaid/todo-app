package middleware

import (
	"strings"

	"todo-app/internal/model"
	"todo-app/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			userUseCase.Log.Warn("Missing Authorization header")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Missing Authorization header",
			})
		}

		// Format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			userUseCase.Log.Warnf("Invalid Authorization format: %s", authHeader)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid Authorization format. Use Bearer <token>",
			})
		}
		token := parts[1]

		// Verifikasi token lewat usecase
		request := &model.VerifyUserRequest{Token: token}
		auth, err := userUseCase.Verify(ctx.UserContext(), request)
		if err != nil {
			userUseCase.Log.Warnf("Failed verify user token: %+v", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid or expired token",
			})
		}

		// Simpan user ke context
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	if v := ctx.Locals("auth"); v != nil {
		return v.(*model.Auth)
	}
	return nil
}
