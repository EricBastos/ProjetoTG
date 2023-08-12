package dtos

import (
	"time"
)

type CreateUserInput struct {
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
	Email     string    `json:"email"`
	TaxId     string    `json:"taxId"`
	CreatedAt time.Time `json:"createdAt"`
}
