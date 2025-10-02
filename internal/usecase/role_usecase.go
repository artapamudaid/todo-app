package usecase

import (
	"context"
	"todo-app/internal/entity"
	"todo-app/internal/model"
	"todo-app/internal/model/converter"
	"todo-app/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	RoleRepository *repository.RoleRepository
}

func NewRoleUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	roleRepository *repository.RoleRepository) *RoleUseCase {
	return &RoleUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		RoleRepository: roleRepository,
	}
}

func (c *RoleUseCase) Create(ctx context.Context, request *model.CreateRoleRequest) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	role := &entity.Role{
		ID:   uuid.New().String(),
		Name: request.Name,
	}

	if err := c.RoleRepository.Create(tx, role); err != nil {
		c.Log.WithError(err).Error("error creating role")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating role")
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}

func (c *RoleUseCase) Update(ctx context.Context, request *model.UpdateRoleRequest) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	c.Log.Debug(request.ID)

	role := new(entity.Role)
	if err := c.RoleRepository.FindById(tx, role, request.ID); err != nil {
		c.Log.WithError(err).Error("error getting role")
		return nil, fiber.ErrNotFound
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	role.Name = request.Name

	if err := c.RoleRepository.Update(tx, role); err != nil {
		c.Log.WithError(err).Error("error updating role")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating role")
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}

func (c *RoleUseCase) Get(ctx context.Context, request *model.GetRoleRequest) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	c.Log.Debug(request.ID)

	role := new(entity.Role)
	if err := c.RoleRepository.FindById(tx, role, request.ID); err != nil {
		c.Log.WithError(err).Error("error getting role")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error getting role")
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}

func (c *RoleUseCase) SoftDelete(ctx context.Context, request *model.GetRoleRequest) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	role := new(entity.Role)
	if err := c.RoleRepository.FindById(tx, role, request.ID); err != nil {
		return nil, fiber.ErrNotFound
	}

	if err := c.RoleRepository.SoftDelete(tx, role); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}

func (c *RoleUseCase) RecycleBin(ctx context.Context, request *model.SearchRoleRequest) ([]model.RoleResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	roles, total, err := c.RoleRepository.SearchTrashed(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting trashed roles")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error committing trashed roles")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = *converter.RoleToResponse(&role)
	}

	return responses, total, nil
}

func (c *RoleUseCase) Restore(ctx context.Context, request *model.GetRoleRequest) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.RoleRepository.Restore(tx, request.ID); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	role := new(entity.Role)
	if err := c.RoleRepository.FindById(tx, role, request.ID); err != nil {
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}

func (c *RoleUseCase) ForceDelete(ctx context.Context, request *model.DeleteRoleRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.RoleRepository.ForceDelete(tx, request.ID); err != nil {
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *RoleUseCase) Search(ctx context.Context, request *model.SearchRoleRequest) ([]model.RoleResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	roles, total, err := c.RoleRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting roles")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error getting roles")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = *converter.RoleToResponse(&role)
	}

	return responses, total, nil
}
