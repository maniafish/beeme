package models

import (
	"beeme/conf"
	"errors"
	"fmt"
	"gas2/pkg/counter"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	rowid        = counter.New()
	userOrmer    orm.Ormer
	UserNotExist = errors.New("User not Exist")
)

// User user model
type User struct {
	ID       int    `json:"id" orm:"pk;column(id);auto"`
	Username string `json:"username"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}

// Init Register database
func Init() error {
	err := orm.RegisterDataBase("default", "mysql", conf.Config.UserDB, conf.Config.DBMaxIdleConns, conf.Config.DBMaxOpenConns)
	if err != nil {
		return fmt.Errorf("RegisterDataBase error: %v", err)
	}

	orm.RegisterModel(new(User))
	orm.RunSyncdb("default", false, true)
	userOrmer = orm.NewOrm()
	return nil
}

// Add add new user
func (u *User) Add() (int, error) {
	_, err := userOrmer.Insert(u)
	if err != nil {
		return -1, err
	}
	return u.ID, nil
}

// Get get user by id
func (u *User) Get() error {
	err := userOrmer.Read(u, "id")
	switch err {
	case orm.ErrNoRows:
		return UserNotExist
	default:
		return err
	}
}

// Update update user
func (u *User) Update() error {
	affect, err := userOrmer.Update(u)
	switch {
	case err != nil:
		return err
	case affect == 0:
		return UserNotExist
	default:
		return nil
	}
}

func (u *User) Delete() error {
	affect, err := userOrmer.Delete(u)
	switch {
	case err != nil:
		return err
	case affect == 0:
		return UserNotExist
	default:
		return nil
	}
}
