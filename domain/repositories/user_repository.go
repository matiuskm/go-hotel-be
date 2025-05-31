package repositories

import "matiuskm/go-hotel-be/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	GetByUsername(username string) (*entities.User, error)
	ExistsByUsername(username string) (bool, error)
	GetByID(id uint) (*entities.User, error)
	GetAll() ([]*entities.User, error)
	Update(user *entities.User) error
	Delete(id uint) error
}