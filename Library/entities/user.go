package entities

import (
	"github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID *entities.ID `json:"id"`

	Email    string `json:"email"`
	Password string `json:"password"`

	TaxId string `json:"taxId"`

	CreatedAt time.Time `json:"createdAt" gorm:"<-:create"`
}

func NewUser(email, password, taxId string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	id := entities.NewID()
	return &User{
		ID:       &id,
		Email:    email,
		TaxId:    taxId,
		Password: string(hash),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
