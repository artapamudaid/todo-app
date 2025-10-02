package http

import (
	"math"
	"todo-app/internal/model"
	"todo-app/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RoleController struct {
	UseCase *usecase.RoleUseCase
	Log     *logrus.Logger
}

func NewRoleController(useCase *usecase.RoleUseCase, log *logrus.Logger) *RoleController {
	return &RoleController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *RoleController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateRoleRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: response})
}

func (c *RoleController) List(ctx *fiber.Ctx) error {

	request := &model.SearchRoleRequest{
		Name: ctx.Query("name", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching role")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.RoleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *RoleController) Get(ctx *fiber.Ctx) error {

	request := &model.GetRoleRequest{
		ID: ctx.Params("roleId"),
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: response})
}

func (c *RoleController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateRoleRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.ID = ctx.Params("roleId")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: response})
}

func (c *RoleController) SoftDelete(ctx *fiber.Ctx) error {
	req := &model.GetRoleRequest{
		ID: ctx.Params("roleId"),
	}

	resp, err := c.UseCase.SoftDelete(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("error soft deleting role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: resp})
}

func (c *RoleController) RecycleBin(ctx *fiber.Ctx) error {
	request := &model.SearchRoleRequest{
		Name: ctx.Query("name", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.RecycleBin(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.RoleResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *RoleController) Restore(ctx *fiber.Ctx) error {
	req := &model.GetRoleRequest{
		ID: ctx.Params("roleId"),
	}

	resp, err := c.UseCase.Restore(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("error restoring role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: resp})
}

func (c *RoleController) ForceDelete(ctx *fiber.Ctx) error {
	req := &model.DeleteRoleRequest{
		ID: ctx.Params("roleId"),
	}

	if err := c.UseCase.ForceDelete(ctx.UserContext(), req); err != nil {
		c.Log.WithError(err).Error("error force deleting role")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
