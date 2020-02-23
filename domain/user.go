package domain

import "time"

type UserID string

func (c UserID) String() string {
	return string(c)
}

type User struct {
	ID           UserID    `json:"id" bson:"_id,omitempty"`
	Login        string    `json:"login" bson:"login"`
	Password     string    `json:"password" bson:"password"`
	Email        string    `json:"email" bson:"email"`
	Phone        string    `json:"phone" bson:"phone"`
	CreationTime time.Time `json:"creation_time" bson:"creation_time"`
}
