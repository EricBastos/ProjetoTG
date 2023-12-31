package dtos

import (
	"time"
)

type CreateUserInput struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	TaxId           string `json:"taxId"`
}

type GetJwtInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetJwtOutput struct {
	AccessToken string `json:"accessToken"`
}

type GetUserOutput struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	TaxId     string    `json:"taxId"`
	CreatedAt time.Time `json:"createdAt"`
}
