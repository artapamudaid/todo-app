package repository

import (
	"todo-app/internal/entity"
	"todo-app/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleRepository struct {
	Repository[entity.Role]
	Log *logrus.Logger
}

func NewRoleRepository(log *logrus.Logger) *RoleRepository {
	return &RoleRepository{
		Log: log,
	}
}

func (r *RoleRepository) FindById(db *gorm.DB, role *entity.Role, id string) error {
	return db.Where("id = ?", id).First(role).Error
}

func (r *RoleRepository) Search(db *gorm.DB, request *model.SearchRoleRequest) ([]entity.Role, int64, error) {
	var roles []entity.Role

	page := request.Page
	if page < 1 {
		page = 1
	}
	size := request.Size
	if size <= 0 {
		size = 10
	}

	query := db.Model(&entity.Role{})

	if err := query.Scopes(r.FilterRole(request)).
		Offset((page - 1) * size).
		Limit(size).
		Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Scopes(r.FilterRole(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *RoleRepository) FilterRole(request *model.SearchRoleRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name := request.Name; name != "" {
			tx = tx.Where("name LIKE ?", "%"+name+"%")
		}
		return tx
	}
}

func (r *RoleRepository) SoftDelete(db *gorm.DB, role *entity.Role) error {
	return db.Delete(role).Error
}

func (r *RoleRepository) Restore(db *gorm.DB, id string) error {
	return db.Unscoped().
		Model(&entity.Role{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

func (r *RoleRepository) ForceDelete(db *gorm.DB, id string) error {
	return db.Unscoped().Where("id = ?", id).Delete(&entity.Role{}).Error
}

func (r *RoleRepository) SearchTrashed(db *gorm.DB, request *model.SearchRoleRequest) ([]entity.Role, int64, error) {
	var roles []entity.Role

	page := request.Page
	if page < 1 {
		page = 1
	}
	size := request.Size
	if size <= 0 {
		size = 10
	}

	query := db.Unscoped().Model(&entity.Role{}).Where("deleted_at IS NOT NULL")

	if err := query.Scopes(r.FilterRole(request)).
		Offset((page - 1) * size).
		Limit(size).
		Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Scopes(r.FilterRole(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}
