package controller

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/gin-scaffold/dao"
	"github.com/joexu01/gin-scaffold/dto"
	"github.com/joexu01/gin-scaffold/lib"
	"github.com/joexu01/gin-scaffold/middleware"
	"github.com/joexu01/gin-scaffold/public"
	"net/http"
	"time"
)

type UserController struct{}

func UserControllerRegister(group *gin.RouterGroup) {
	user := &UserController{}
	group.POST("/login", user.UserLogin)
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
// @Router       /user/login [post]
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

	//c.SetCookie("user_token", "123", 3600, "/", "localhost", false, true)

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
	//sessKey := "user_id_" + strconv.Itoa(user.Id)
	session.Set(public.UserSessionKey, string(bytes))
	_ = session.Save()

	out := dto.UserLoginOutput{Token: string(bytes)}
	middleware.ResponseSuccess(c, out)
}
