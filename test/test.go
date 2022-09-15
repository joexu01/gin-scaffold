package main

import (
	"fmt"
	"github.com/joexu01/gin-scaffold/dao"
	"github.com/joexu01/gin-scaffold/dto"
	"github.com/joexu01/gin-scaffold/public"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {
	connStr := "root:atk_2018@tcp(127.0.0.1:3306)/dev_test?charset=utf8mb4&parseTime=true&loc=Asia%2FChongqing"
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	userInsert := dto.NewUserInput{
		Username:    "joexu01",
		RawPassword: "12345678",
		Email:       "joexu01@yahoo.com",
	}

	pwdHash, err := public.GeneratePwdHash([]byte(userInsert.RawPassword))
	if err != nil {
		panic(err)
	}

	user := dao.User{
		Id:        0,
		Username:  userInsert.Username,
		Password:  pwdHash,
		Email:     userInsert.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDelete:  0,
		UserRole:  1,
		Role:      dao.Role{},
	}

	result := db.Create(&user)

	fmt.Println(user.Id)

	fmt.Printf("%+v", result)

	user1 := dao.User{
		Id:        0,
		Username:  "joexu01",
		Password:  "",
		Email:     "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		IsDelete:  0,
	}

	result = db.First(&user1)

	fmt.Println(user1.Email)
	fmt.Printf("%+v", result)

	fmt.Printf("%s", string("hello\n\n\n"))
}
