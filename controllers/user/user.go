package user

import (
	"beeme/models"
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Controller Operations about Users
type Controller struct {
	beego.Controller
}

// Post post
// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {json} {"uid": models.User.ID}
// @Failure 400 invalid body
// @Failure 500 internal error
// @router / [post]
func (u *Controller) Post() {
	logPrefix := "user.Post()"
	var user *models.User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		logs.Error("%s.RequestBody err: %v", logPrefix, err)
		u.Response(400, err.Error())
		return
	}

	uid, err := user.Add()
	if err != nil {
		logs.Error("%s.Add err: %v", logPrefix, err)
		u.Response(500, err.Error())
		return
	}

	u.Response(200, map[string]int{"uid": uid})
}

// Get get
// @Title Get
// @Description get user by uid
// @Param	uid		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 400 invalid uid
// @Failure 404 user not found
// @Failure 500 internal error
// @router /:uid [get]
func (u *Controller) Get() {
	logPrefix := "user.Get()"
	uid, err := u.GetInt(":uid")
	if err != nil {
		logs.Info("%s.GetInt err: %v", logPrefix, err)
		u.Response(400, err.Error())
		return
	}

	user := &models.User{ID: uid}
	err = user.Get()
	switch err {
	case models.ErrUNE:
		u.Response(404, err.Error())
		return
	case nil:
		u.Response(200, user)
		return
	default:
		logs.Error("%s.Get err: %v", logPrefix, err)
		u.Response(500, err.Error())
		return
	}
}

// Put update
// @Title Update
// @Description update the user
// @Param	uid		path 	int	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 400 invalid uid/body
// @Failure 404 user not exist
// @Failure 500 internal error
// @router /:uid [put]
func (u *Controller) Put() {
	logPrefix := "user.Put()"
	uid, err := u.GetInt(":uid")
	if err != nil {
		logs.Info("%s.GetInt err: %v", logPrefix, err)
		u.Response(400, err.Error())
		return
	}

	user := &models.User{ID: uid}
	err = user.Get()
	switch err {
	case models.ErrUNE:
		u.Response(404, err.Error())
		return
	case nil:
		err = json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		if err != nil {
			logs.Info("%s.RequestBody err: %v", logPrefix, err)
			u.Response(400, err.Error())
			return
		}

		user.ID = uid
		err = user.Update()
		switch err {
		case nil:
			u.Response(200, user)
			return
		default:
			logs.Error("%s.Update err: %v", logPrefix, err)
			u.Response(500, err.Error())
			return
		}
	default:
		logs.Error("%s.Get err: %v", logPrefix, err)
		u.Response(500, err.Error())
		return
	}
}

// Delete delete
// @Title Delete
// @Description delete the user
// @Param	uid		path 	int	true		"The uid you want to delete"
// @Success 200 {json} {"uid": models.User.ID}
// @Failure 400 invalid uid
// @Failure 404 user not exist
// @Failure 500 internal error
// @router /:uid [delete]
func (u *Controller) Delete() {
	logPrefix := "user.Delete()"
	uid, err := u.GetInt(":uid")
	user := &models.User{ID: uid}
	err = user.Delete()
	switch err {
	case models.ErrUNE:
		u.Response(404, err.Error())
		return
	case nil:
		u.Response(200, map[string]int{"uid": uid})
		return
	default:
		logs.Error("%s.Delete err: %v", logPrefix, err)
		u.Response(500, err.Error())
		return
	}
}

// Response return
func (u *Controller) Response(code int, ret interface{}) {
	u.Ctx.Output.Status = code
	u.Data["json"] = ret
	u.ServeJSON()
}
