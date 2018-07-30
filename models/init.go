package models

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var (
	ormer orm.Ormer
	// ErrUNE  user not exist error
	ErrUNE = errors.New("User not Exist")
)

// Init 初始化操作
func Init() {
	InitDB()
	InitDir()
}

// InitDB Register database
func InitDB() {
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
		panic(fmt.Sprintf("register database userdb err: %v", err))
	}

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(fmt.Sprintf("run sync db err: %v", err))
	}

	ormer = orm.NewOrm()
}

// InitDir 初始化目录
func InitDir() {
	d := beego.AppConfig.String("apps::FileDir")
	if !strings.HasSuffix(d, "/") {
		panic(fmt.Sprintf("Invalid FileDir: %v", d))
	}

	if strings.HasPrefix(d, "./") {
		_, thisFilePath, _, _ := runtime.Caller(0)
		prefix := filepath.Dir(filepath.Dir(thisFilePath))
		d = filepath.Join(prefix, d[2:])
	}

	_, err := os.Stat(d)
	if err != nil {
		panic(fmt.Sprintf("check file dir err: %v", err))
	}

	beego.AppConfig.Set("apps::FileDir", d)
}
