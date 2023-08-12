package database

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"gorm.io/gorm"
	"sync"
)

type UserDB struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{
		DB: db,
	}
}

func (u *UserDB) Create(user *entities.User) error {
	var userByEmail *entities.User
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		userByEmail, _ = u.FindByEmail(user.Email)
		wg.Done()
	}()
	wg.Wait()
	foundUser := userByEmail

	if foundUser == nil {
		return u.DB.Create(user).Error
	} else {
		return errors.New(utils.UserAlreadyExists)
	}
}

func (u *UserDB) FindByEmail(email string) (*entities.User, error) {
	var user entities.User

	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	} else if user.ID == nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (u *UserDB) FindById(id string) (*entities.User, error) {
	var user entities.User

	if err := u.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	} else if user.ID == nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
