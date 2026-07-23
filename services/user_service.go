package services

import (
	"errors"

	"github.com/Yohannes3003/project-management2/models"
	"github.com/Yohannes3003/project-management2/repositories"
	"github.com/Yohannes3003/project-management2/utils"
	"github.com/google/uuid"
)

type UserService interface {
	Register(user *models.User) error
	Login (email, password string) (*models.User, error)
	GetById (id uint) (*models.User, error)
	GetByPublicId (id string) (*models.User, error)
	GetAllPagination (filter, sort string, limit, offset int)([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService)Register(user *models.User) error {

	existingUser , _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}
	hased, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hased
	user.Role = "user"
	user.PublicID = uuid.New()
	return s.repo.Create(user)

}

func (s *userService) Login (email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credential")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credential")
	}
	return user, nil
}

func (s *userService) GetById (id uint) (*models.User, error) {
	return s.repo.FindById(id)
}

func (s *userService) GetByPublicId (id string) (*models.User, error) {
	return s.repo.FindByPublicId(id)
}

func (s *userService) GetAllPagination (filter, sort string, limit, offset int)([]models.User, int64, error) {
	return s.repo.FindAllPagination(filter, sort, limit, offset)
}

func (s *userService) Update(user *models.User) error {
	return s.repo.Update(user)
}

func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
}	