package userRepository

import (
	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *GormUserRepository) FindUserByEmail(users []models.User, email string) *models.User {
	for _, u := range users {
		if u.Email == email {
			return &u
		}
	}
	return nil
}

func (r *GormUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) SaveUser(newUser models.User) error {
	return r.db.Create(&newUser).Error
}

func (r *GormUserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) SaveAllUsers(users []models.User) error {
	for _, user := range users {
		if err := r.db.Save(&user).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *GormUserRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) ToggleStaffAvailability(userID string) error {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}
	user.Available = !user.Available
	return r.db.Save(&user).Error
}

func (r *GormUserRepository) GetStaffAvailability(userID string) (bool, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return false, err
	}
	return user.Available, nil
}
