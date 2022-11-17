package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string            `json:"Password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token"`
	User_type     string             `json:"user_type"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    int64              `json:"created_at"`
	Updated_at    int64              `json:"updated_at"`
	User_id       string             `json:"user_id"`
}

type CreateUser struct {
	First_name *string `json:"first_name" validate:"required,min=2,max=100"`
	Last_name  *string `json:"last_name" validate:"required,min=2,max=100"`
	Password   *string `json:"Password" validate:"required,min=6"`
	Email      *string `json:"email" validate:"email,required"`
	Phone      *string `json:"phone" validate:"required"`
}

type UserRes struct {
	Token         *string `json:"token"`
	Refresh_token *string `json:"refresh_token"`
}

type UserLogin struct {
	Email    *string `json:"email" validate:"email,required"`
	Password *string `json:"Password" validate:"required,min=6"`
}

type UserUpdate struct {
	First_name *string `json:"first_name"`
	Last_name  *string `json:"last_name"`
	Password   *string `json:"Password"`
	Email      *string `json:"email" validate:"email"`
	Phone      *string `json:"phone"`
}
