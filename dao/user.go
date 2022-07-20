package dao

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/gin-scaffold/dto"
	"github.com/joexu01/gin-scaffold/public"
	"gorm.io/gorm"
	"time"
)

/*

CREATE TABLE `user` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
    `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
    `hashed_password` varchar(512) NOT NULL DEFAULT '' COMMENT '加盐后密码',
    `created_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
    `updated_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
    `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

*/

type User struct {
	Id        int       `json:"id" gorm:"column:id"`
	Username  string    `json:"username" gorm:"column:username"`
	Password  string    `json:"password" gorm:"column:hashed_password"`
	Email     string    `json:"email" gorm:"column:email"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) PwdCheck(rawPwd string) bool {
	return public.ComparePwdAndHash([]byte(rawPwd), u.Password)
}

// LoginCheck 对用户输入数据进行检查  决定是否允许登陆
func (u *User) LoginCheck(c *gin.Context, dbConn *gorm.DB, param *dto.UserLoginInput) (*User, error) {
	userInfo, err := u.Find(c, dbConn, &User{Username: param.Username, IsDelete: 0})
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}

	if !userInfo.PwdCheck(param.Password) {
		return nil, errors.New("密码错误")
	}

	return userInfo, nil
}

func (u *User) Find(_ *gin.Context, tx *gorm.DB, search *User) (*User, error) {
	out := &User{}
	err := tx.Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (u *User) Save(_ *gin.Context, tx *gorm.DB) error {
	return tx.Save(u).Error
}
