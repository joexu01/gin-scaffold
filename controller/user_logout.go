package controller

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/gin-scaffold/middleware"
	"github.com/joexu01/gin-scaffold/public"
)

type UserLogoutController struct{}

func UserLogoutRegister(group *gin.RouterGroup) {
	user := &UserLogoutController{}
	group.GET("/logout", user.UserLogout)
}

// UserLogout godoc
// @Summary      用户登出
// @Description  就是用户登出呗
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  middleware.Response
// @Failure      500  {object}  middleware.Response
// @Router       /user-logout/logout [get]
func (u *UserLogoutController) UserLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(public.UserSessionKey)
	_ = session.Save()
	middleware.ResponseSuccess(c, "")
}
