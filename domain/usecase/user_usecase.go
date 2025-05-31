package usecase

import "matiuskm/go-hotel-be/domain/entities"

type UserUsecase interface {
	AddUser(username, password, fullName, role string) error
	EditUser(id uint, fullName, role string) error
	DeleteUser(id uint) error
	AssignRole(id uint, role string) error
	Login(username, password string) (*entities.User, error)
}