package dto

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/joexu01/gin-scaffold/lib"
)

type DefaultClaims struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (c *DefaultClaims) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(lib.GetStringConf("base.jwt.jwt_secret")))
	if err != nil {
		return "", err
	}
	return ss, nil
}
