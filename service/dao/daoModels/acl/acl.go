package aclDaoModel

import (
	"time"
)

type FieldName string

func (f FieldName) String() string {
	return string(f)
}

const (
	Username     FieldName = "username"
	Password     FieldName = "password"
	MobileNumber FieldName = "mobile_number"
	Email        FieldName = "email"
	Sex          FieldName = "sex"
	Birthday     FieldName = "birthday"
	CreatedAt    FieldName = "created_at"
)

type User struct {
	Username     string    `bson:"username" json:"username"`
	Password     string    `bson:"password,omitempty" json:"password"`
	MobileNumber string    `bson:"mobile_number" json:"mobile_number"`
	Email        string    `bson:"email" json:"email"`
	Sex          string    `bson:"sex" json:"sex"`
	Birthday     string    `bson:"birthday" json:"birthday"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	Token        string    `bson:"token" json:"token"`
}

type Query struct {
	BulkUserArgs []BulkUserArg
	CreatedAt    time.Time
}

type BulkUserArg struct {
	Username     string `json:"username"`
	Password     string `bson:"password,omitempty" json:"password"`
	MobileNumber string `bson:"mobile_number" json:"mobile_number"`
	Email        string `bson:"email" json:"email"`
	Sex          string `bson:"sex" json:"sex"`
	Birthday     string `bson:"birthday" json:"birthday"`
}

type UserSession struct {
	Username string `bson:"username"`
	Token    string `bson:"token"`
}

type UserWithMeta struct {
	BulkUserArg `bson:",inline"`
	CreatedAt   time.Time `bson:"created_at"`
}
