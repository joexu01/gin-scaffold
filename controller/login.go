package controller

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joexu01/gin-scaffold/dao"
	"github.com/joexu01/gin-scaffold/dto"
	"github.com/joexu01/gin-scaffold/lib"
	"github.com/joexu01/gin-scaffold/middleware"
	"github.com/joexu01/gin-scaffold/public"
	"net/http"
	"time"
)

type LoginController struct{}

func LoginControllerRegister(group *gin.RouterGroup) {
	user := &LoginController{}
	group.POST("/login", user.UserLogin)
	//group.POST("/register", user.)
}

// UserLogin godoc
// @Summary      用户登录
// @Description  就是用户登录呗
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        struct body dto.UserLoginInput true "用户登录输入"
// @Success      200  {object}  middleware.Response{data=dto.UserLoginOutput} "success"
// @Failure      500  {object}  middleware.Response
// @Router       /login [post]
func (l *LoginController) UserLogin(c *gin.Context) {
	params := new(dto.UserLoginInput)
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseWithCode(c, http.StatusInternalServerError, 2001, err, "")
		return
	}

	user := &dao.User{}
	user, err = user.LoginCheck(c, db, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	claims := &dto.DefaultClaims{
		Id:       user.Id,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "self",
			Subject:   "login",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 2)),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "",
		},
	}

	token, err := claims.GenerateToken()
	if err != nil {
		middleware.ResponseError(c, 2500, err)
		return
	}

	output := &dto.UserLoginOutput{Token: token}

	session := sessions.Default(c)
	session.Set(public.UserSessionKey, token)
	_ = session.Save()

	c.Header("Authorization", "Bearer "+token)
	middleware.ResponseSuccess(c, output)
}
