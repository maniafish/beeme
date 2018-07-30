package models

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	ormer orm.Ormer
	// ErrUNE  user not exist error
	ErrUNE = errors.New("User not Exist")
)

// Init Register database
func Init() {
	maxIdleConns, _ := beego.AppConfig.Int("apps::MaxIdleConns")
	if maxIdleConns <= 0 {
		maxIdleConns = 16
	}

	maxOpenConns, _ := beego.AppConfig.Int("apps::DBMaxOpenConns")
	if maxOpenConns <= 0 {
		maxOpenConns = 16
	}

	orm.RegisterModel(new(DemoUser))
	err := orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("apps::UserDB"), maxIdleConns, maxOpenConns)
	if err != nil {
		panic(err)
	}

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}

	ormer = orm.NewOrm()
}

// DemoUser DemoUser model
type DemoUser struct {
	ID       int    `json:"id" orm:"pk;column(id);auto"`
	Username string `json:"username"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}

// Add add new DemoUser
func (u *DemoUser) Add() (int, error) {
	_, err := ormer.Insert(u)
	if err != nil {
		return -1, err
	}
	return u.ID, nil
}

// Get get DemoUser by id
func (u *DemoUser) Get() error {
	err := ormer.Read(u, "id")
	switch err {
	case orm.ErrNoRows:
		return ErrUNE
	default:
		return err
	}
}

// Update update DemoUser
func (u *DemoUser) Update() error {
	_, err := ormer.Update(u)
	return err
}

// Delete delete DemoUser
func (u *DemoUser) Delete() error {
	affect, err := ormer.Delete(u)
	switch {
	case err != nil:
		return err
	case affect == 0:
		return ErrUNE
	default:
		return nil
	}
}
