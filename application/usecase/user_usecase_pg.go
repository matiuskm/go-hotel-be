package usecase

import (
	"errors"
	"matiuskm/go-hotel-be/domain/entities"
	"matiuskm/go-hotel-be/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecasePG struct {
	UserRepo repositories.UserRepository
}

func (u *UserUsecasePG) AddUser(username, password, fullName, role string) error {
	exist, err := u.UserRepo.ExistsByUsername(username)
	if err!= nil {
		return errors.New("database error")
	}
	if exist {
		return errors.New("username already exists")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &entities.User{
		Username: username,
		Password: string(hashed),
		FullName: fullName,
		Role:     role,
	}
	return u.UserRepo.Create(user)
}

func (u *UserUsecasePG) EditUser(id uint, fullName, role string) error {
	user, err := u.UserRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	user.FullName = fullName
	user.Role = role
	return u.UserRepo.Update(user)
}

func (u *UserUsecasePG) DeleteUser(id uint) error {
	return u.UserRepo.Delete(id)
}

func (u *UserUsecasePG) AssignRole(id uint, role string) error {
	user, err := u.UserRepo.GetByID(id)
	if err!= nil {
		return errors.New("user not found")
	}
	user.Role = role
	return u.UserRepo.Update(user)
}

func (u *UserUsecasePG) Login(username, password string) (*entities.User, error) {
	user, err := u.UserRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}