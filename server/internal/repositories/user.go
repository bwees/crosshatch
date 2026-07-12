package repositories

import (
	"crosshatch/internal/database/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	users := []models.User{}
	err := r.db.Order("username").Find(&users).Error
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (r *UserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := models.User{}
	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	user := models.User{}
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
