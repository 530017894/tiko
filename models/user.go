package models

import (
	"errors"
	"fmt"
	"myapp/libs"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null; type:varchar(60)" json:"name" validate:"required,gte=2,lte=50" comment:"用户名"`
	Username string `gorm:"unique;not null;type:varchar(60)" json:"username" validate:"required,gte=2,lte=50"  comment:"名称"`
	Password string `gorm:"type:varchar(100)" json:"password"  comment:"密码"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"introduction" comment:"简介"`
	Avatar   string `gorm:"type:longText" json:"avatar"  comment:"头像"`
}

func NewUser() *User {
	return &User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// GetUser get user
func GetUser(search *Search) (*User, error) {
	t := NewUser()
	err := Found(search).First(t).Error
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// DeleteUser del user . if user's username is username ,can't del it.
func DeleteUser(id uint) error {
	s := &Search{
		Fields: []*Filed{
			{
				Key:       "id",
				Condition: "=",
				Value:     id,
			},
		},
	}
	u, err := GetUser(s)
	if err != nil {
		return err
	}
	if u.Username == "username" {
		return errors.New(fmt.Sprintf("不能删除管理员 : %s \n ", u.Username))
	}

	if err := libs.Db.Delete(u, id).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteUserByIdErr:%s \n ", err))
		return err
	}
	return nil
}

// GetAllUsers get all users
func GetAllUsers(s *Search) ([]*User, int64, error) {
	var users []*User
	var count int64
	q := GetAll(&User{}, s)
	if err := q.Count(&count).Error; err != nil {
		return nil, count, err
	}
	q = q.Scopes(Paginate(s.Offset, s.Limit), Relation(s.Relations))
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllUserErr:%s \n ", err))
		return nil, count, err
	}
	return users, count, nil
}

// CreateUser create user
func (u *User) CreateUser() error {
	u.Password = libs.HashPassword(u.Password)
	if err := libs.Db.Create(u).Error; err != nil {
		return err
	}

	return nil
}

// UpdateUserById update user by id
func UpdateUserById(id uint, nu *User) error {
	if len(nu.Password) > 0 {
		nu.Password = libs.HashPassword(nu.Password)
	}
	if err := Update(&User{}, nu, id); err != nil {
		return err
	}
	return nil
}
