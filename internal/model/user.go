package model

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidPassword = errors.New("invalid username or password")
	ErrTokenInvalid    = errors.New("token invalid")
	ErrUserNotFound    = errors.New("user not found")
)

type User struct {
	Id         string `xorm:"pk"`
	Name       string `xorm:"unique"`
	Salt       string
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastSeenAt time.Time
}

type UserClaims struct {
	UserId string
	jwt.MapClaims
}

// ----------------------
// QUERIES

type GetUserByNameQuery struct {
	Ctx *ReqContext

	Name string

	User *User
}

type LoginQuery struct {
	Ctx *ReqContext

	Name     string
	Password string

	User *User
}

type GetUserByIdQuery struct {
	Ctx *ReqContext

	UserId string

	User *User
}

type GetSignedInUserQuery struct {
	Ctx *ReqContext

	UserId string

	User *SignedInUser
}

// ----------------------
// COMMANDS

type CreateUserCommand struct {
	Ctx *ReqContext

	User *User
}

type SignUpCommand struct {
	Ctx *ReqContext

	Name     string
	Password string

	User *User
}

type UpdateUserLastSeenAtCommand struct {
	Ctx *ReqContext

	UserId string
}

type CreateTokenCommand struct {
	User *User

	AccessToken  string
	RefreshToken string
}

type RefreshTokenCommand struct {
	RefreshToken string
	AccessToken  string
}

type LookupTokenCommand struct {
	Token string

	Claims *UserClaims
}

// ----------------------
// DTO

type SignedInUser struct {
	UserId     string
	Name       string
	LastSeenAt time.Time
}

func (u *SignedInUser) ShouldUpdateLastSeenAt() bool {
	return time.Since(u.LastSeenAt) > 30*time.Minute
}
