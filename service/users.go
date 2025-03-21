package service

import (
	"errors"
	"ecommerce-project-go/entity"
	"ecommerce-project-go/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(input entity.InputRegisterUsers) (entity.Users, error)
	Login(input entity.InputLogin) (entity.Users, error)
	UpdateUser(id int, input entity.InputUpdateUser) (entity.Users, error)
	GetUserById(id int) (entity.Users, error)
	DeleteUser(id int) error
	GetAll(isAdmin bool, page int, limit int) ([]entity.Users, map[string]interface{}, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{userRepository}
}

func (s *userService) RegisterUser(input entity.InputRegisterUsers) (entity.Users, error) {
	var user entity.Users

	user.FullName = input.FullName
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)

	// checking if email already exist
	_, emailExist, _ := s.userRepository.FindByEmail(user.Email)

	// if email not available or email already exist
	if emailExist {
		return user, errors.New("email already exist")
	}

	// if email available
	newUser, err := s.userRepository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *userService) Login(input entity.InputLogin) (entity.Users, error) {
	email := input.Email
	pwd := input.Password

	// check if the email exist
	user, _, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return user, errors.New("user not found")
	}

	if user.ID == 0 {
		return user, errors.New("user not found")
	}

	// compare user input password with password hash in database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pwd))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) UpdateUser(id int, input entity.InputUpdateUser) (entity.Users, error) {
	user, err := s.userRepository.FindById(id)
	if err != nil {
		return user, errors.New("user with id not found")
	}

	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Email != "" {
		// email validation if email exist
		existingUser, emailExist, _ := s.userRepository.FindByEmail(input.Email)

		if emailExist && existingUser.ID != user.ID {
			return user, errors.New("email already exists")
		}

		user.Email = input.Email
	}
	if input.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
		if err != nil {
			return user, err
		}
		user.PasswordHash = string(passwordHash)
	}

	if input.IsAdmin != user.IsAdmin {
		user.IsAdmin = input.IsAdmin
	}

	updatedUser, err := s.userRepository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *userService) GetUserById(id int) (entity.Users, error) {
	user, err := s.userRepository.FindById(id)
	if err != nil {
		return user, errors.New("user not found with that id")
	}

	if user.ID == 0 {
		return user, errors.New("user not found with that id")
	}

	return user, nil
}

func (s *userService) DeleteUser(id int) error {
	user, err := s.userRepository.FindById(id)
	if err != nil {
		return errors.New("user with id not found")
	}

	err = s.userRepository.Delete(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetAll(isAdmin bool, page int, limit int) ([]entity.Users, map[string]interface{}, error) {
    if !isAdmin {
        return nil, nil, errors.New("you're not authorized")
    }

    users, meta, err := s.userRepository.GetAll(page, limit)
    if err != nil {
        return users, nil, err
    }

    return users, meta, nil
}

