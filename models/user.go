package models

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	userOrmer orm.Ormer
	// ErrUNE  user not exist error
	ErrUNE = errors.New("User not Exist")
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
func Init() {
	maxIdleConns, _ := beego.AppConfig.Int("apps::MaxIdleConns")
	if maxIdleConns <= 0 {
		maxIdleConns = 16
	}

	maxOpenConns, _ := beego.AppConfig.Int("DBMaxOpenConns")
	if maxOpenConns <= 0 {
		maxOpenConns = 16
	}

	err := orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("apps::UserDB"), maxIdleConns, maxOpenConns)
	if err != nil {
		panic(err)
	}

	orm.RegisterModel(new(User))
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
	userOrmer = orm.NewOrm()
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
		return ErrUNE
	default:
		return err
	}
}

// Update update user
func (u *User) Update() error {
	_, err := userOrmer.Update(u)
	return err
}

// Delete delete user
func (u *User) Delete() error {
	affect, err := userOrmer.Delete(u)
	switch {
	case err != nil:
		return err
	case affect == 0:
		return ErrUNE
	default:
		return nil
	}
}
