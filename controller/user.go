package controller

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/gin-scaffold/dao"
	"github.com/joexu01/gin-scaffold/dto"
	"github.com/joexu01/gin-scaffold/lib"
	"github.com/joexu01/gin-scaffold/middleware"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
}

func UserControllerRegister(group *gin.RouterGroup) {
	user := &UserController{}
	group.POST("/login", user.UserLogin)
}

func (u *UserController) UserLogin(c *gin.Context) {
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

	c.SetCookie("user_token", "123", 3600, "/", "localhost", false, true)

	sessInfo := &dto.UserSessionInfo{
		Id:        user.Id,
		UserName:  user.Username,
		LoginTime: time.Now(),
	}

	bytes, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	session := sessions.Default(c)
	session.Set("user_id_"+strconv.Itoa(user.Id), string(bytes))
}
