package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/joexu01/gin-scaffold/public"
	"time"
)

type UserInfoOutput struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

type ChangePwdInput struct {
	OldPassword string `json:"old_password" form:"old_password" comment:"旧密码" example:"123456" validate:"required"` // 旧密码
	NewPassword string `json:"new_password" form:"new_password" comment:"新密码" example:"123456" validate:"required"` // 新密码
}

func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.GetValidParamsDefault(c, param)
}

type UserLoginInput struct {
	Username string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required,validate_username"` // 管理员用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                   // 管理员密码
}

type UserSessionInfo struct {
	Id        int       `json:"id"`
	UserName  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}

func (param *UserLoginInput) BindValidParam(c *gin.Context) error {
	return public.GetValidParamsDefault(c, param)
}

type UserLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` // 返回的Token
}

type NewUserInput struct {
	Username    string `json:"username" validate:"required"`
	RawPassword string `json:"raw_password" validate:"required"`
	Email       string `json:"email" validate:"required"`
}

func (param *NewUserInput) BindValidParam(c *gin.Context) error {
	return public.GetValidParamsDefault(c, param)
}
