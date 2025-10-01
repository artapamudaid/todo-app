package usecase

import (
	"context"
	"todo-app/internal/entity"
	"todo-app/internal/gateway/messaging"
	"todo-app/internal/model"
	"todo-app/internal/model/converter"
	"todo-app/internal/repository"
	"todo-app/internal/util/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	UserProducer   *messaging.UserProducer
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, userProducer *messaging.UserProducer) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		UserProducer:   userProducer,
	}
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("Failed find user by token : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.Auth{ID: user.ID}, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Generate new UUID for user ID
	userId := uuid.New().String()

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Check if email already exists
	existingUser := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, existingUser, request.Email); err == nil {
		c.Log.Warnf("Email already exists: %s", request.Email)
		return nil, fiber.ErrConflict
	}

	user := &entity.User{
		ID:       userId,
		Email:    request.Email,
		Password: string(password),
		Name:     request.Name,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user created event")
		if err = c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user created event : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user created event")
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validasi request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body  : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Cari user berdasarkan email
	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	// Bandingkan password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Invalid password : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	// Generate Access Token (JWT 24 jam)
	accessToken, err := helper.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.Log.Warnf("Failed to generate JWT access token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate Refresh Token (UUID/random string)
	user.Token = accessToken

	// Simpan refresh token ke DB
	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Publish event kalau Kafka aktif
	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user login event")
		if err := c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user login event : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user login event")
	}

	expiresIn := int64(helper.JWT_EXPIRATION_HOURS * 3600)

	// Return user + token ke client
	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Token:     accessToken,
		ExpiresIn: expiresIn, // dalam detik
	}, nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Refresh(ctx context.Context, refreshToken string) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if refreshToken == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "refresh token required")
	}

	// Cari user berdasarkan refresh token
	user := new(entity.User)
	if err := c.UserRepository.FindByToken(tx, user, refreshToken); err != nil {
		c.Log.Warnf("Refresh failed, user not found with token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	// Generate access token baru (expire 24 jam)
	accessToken, err := helper.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.Log.Errorf("Failed generate new JWT : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user.Token = accessToken // ganti field di entity jadi RefreshToken

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Errorf("Failed update user refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Publish event ke Kafka (opsional)
	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user refresh token event")
		if err := c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user refresh token event : %+v", err)
		}
	}

	expiresIn := int64(helper.JWT_EXPIRATION_HOURS * 3600)

	// Return token baru
	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Token:     accessToken,
		ExpiresIn: expiresIn,
	}, nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return false, fiber.ErrNotFound
	}

	user.Token = ""

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user logout event")
		if err := c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user logout event : %+v", err)
			return false, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user logout event")
	}

	return true, nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Check if email already exists (only if email is being updated)
	if request.Email != "" {
		existingUser := new(entity.User)
		if err := c.UserRepository.FindByEmail(tx, existingUser, request.Email); err == nil {
			if existingUser.ID != user.ID {
				c.Log.Warnf("Email already exists: %s", request.Email)
				return nil, fiber.ErrConflict
			}
		}
		user.Email = request.Email
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(password)
	}

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user updated event")
		if err := c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user updated event : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user updated event")
	}

	return converter.UserToResponse(user), nil
}
