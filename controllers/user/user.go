package user

import (
	"beeme/models"
	"encoding/json"

	"github.com/astaxie/beego"
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
	var user *models.User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		u.Response(400, err.Error())
	}

	uid, err := user.Add()
	if err != nil {
		u.Response(500, err.Error())
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
	uid, err := u.GetInt(":uid")
	if err != nil {
		u.Response(400, err.Error())
	}

	user := &models.User{ID: uid}
	err = user.Get()
	switch err {
	case models.ErrUNE:
		u.Response(404, err.Error())
	case nil:
		u.Response(200, user)
	default:
		u.Response(500, err.Error())
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
	uid, err := u.GetInt(":uid")
	if err != nil {
		u.Response(400, err.Error())
	}

	user := &models.User{ID: uid}
	err = user.Get()
	switch err {
	case models.ErrUNE:
		u.Response(404, err.Error())
	case nil:
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		user.ID = uid
		err = user.Update()
		switch err {
		case nil:
			u.Response(200, user)
		default:
			u.Response(500, err.Error())
		}
	default:
		u.Response(500, err.Error())
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
	uid, err := u.GetInt(":uid")
	user := &models.User{ID: uid}
	err = user.Delete()
	switch err {
	case models.ErrUNE:
		u.Response(404, err.Error())
	case nil:
		u.Response(200, map[string]int{"uid": uid})
	default:
		u.Response(500, err.Error())
	}
}

// Response return
func (u *Controller) Response(code int, ret interface{}) {
	u.Ctx.Output.Status = code
	u.Data["json"] = ret
	u.ServeJSON()
}
