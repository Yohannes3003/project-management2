package repositories

import (
	"strings"

	"github.com/Yohannes3003/project-management2/config"
	"github.com/Yohannes3003/project-management2/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindById (id uint) (*models.User, error)
	FindByPublicId (publicID string) (*models.User, error)
	FindAllPagination (filter, sort string, limit, offset int)([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct {

}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create (user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) FindByEmail (email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ? ", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindById (id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPublicId (publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id = ?", publicID).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAllPagination (filter, sort string, limit, offset int)([]models.User, int64, error) {
	var users []models.User
	var total int64
	
	db := config.DB.Model(&models.User{})

	//filtering
	if filter != "" {
		filterPattern := "%" + filter + "%" 
		db = db.Where("name i like ? Or email i like ?", filterPattern, filterPattern)
	}
	//counting data
	if err := db.Count(&total).Error ; err != nil {
		return nil, 0,err
	}
	//sorting
	if sort != "" {
		//Misalnya sort = name (ascending) sort = -name (descending)
		if sort == "-id"{
			sort = "-internal_id"
		} else if sort == "id" {
			sort = "internal_id"
		}

		if strings.HasPrefix(sort, "-") {
			sort = strings.TrimPrefix(sort, "-") + " DESC "
		} else {
			sort += " ASC "
		}

		db = db.Order(sort)
	}

	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err 
}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Model(&models.User{}).Where("public_id = ?", user.PublicID).Updates(map[string]interface{}{
		"name":  user.Name,
	}).Error
}

func (r *userRepository) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}