package database

import (
	"matiuskm/go-hotel-be/domain/entities"

	"gorm.io/gorm"
)

type UserRepoPG struct {
	DB *gorm.DB
}

func (r *UserRepoPG) Create(user *entities.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepoPG) GetByUsername(username string) (*entities.User, error) {
	var user entities.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepoPG) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.User{}).Where("username =?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepoPG) GetByID(id uint) (*entities.User, error) {
	var user entities.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepoPG) GetAll() ([]*entities.User, error) {
	var users []*entities.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepoPG) Update(user *entities.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepoPG) Delete(id uint) error {
	return r.DB.Delete(&entities.User{}, id).Error
}