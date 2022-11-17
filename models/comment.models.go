package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentAuthor struct {
	User_id string `json:"user_id"`
	Email   string `json:"email"`
}

type CommentModels struct {
	Answer string `json:"answer" validate:"required"`
}
type Comment struct {
	ID         primitive.ObjectID `bson:"_id"`
	Answer     string             `json:"answer"`
	Created_at int64              `json:"created_at"`
	Updated_at int64              `json:"updated_At"`
	Note_id    string             `json:"note_id"`
	Author     CommentAuthor      `json:"author"`
}
