package middleware

import (
	"errors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joexu01/gin-scaffold/dto"
	"github.com/joexu01/gin-scaffold/lib"
	"github.com/joexu01/gin-scaffold/public"
	"log"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		tokenStr := string(session.Get(public.UserSessionKey).(string))
		log.Println(tokenStr)

		token, err := jwt.ParseWithClaims(tokenStr, &dto.DefaultClaims{}, func(token *jwt.Token) (interface{}, error) { // 解析token
			return []byte(lib.GetStringConf("base.jwt.jwt_secret")), nil
		})
		if err != nil {
			ResponseError(c, 2400, err)
			c.Abort()
		}
		if claims, ok := token.Claims.(*dto.DefaultClaims); ok && token.Valid { // 校验token
			c.Set(public.ContextKeyUserId, claims.Id)
			c.Next()
		} else {
			ResponseError(c, 2400, errors.New("invalid token"))
			c.Abort()
		}
	}
}
