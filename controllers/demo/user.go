package demo

import (
	"beeme/controllers"
	"beeme/models"
	"encoding/json"
)

// UserController Operations about Users
type UserController struct {
	controllers.Controller
}

// Prepare set controller param
func (u *UserController) Prepare() {
	u.Name = "user_operator"
	u.Controller.Prepare()
}

// Post post
// @Title CreateUser
// @Description create users
// @Param	body		body 	models.DemoUser	true		"body for user content"
// @Success 200 {json} {"uid": models.DemoUser.ID}
// @Failure 400 invalid body
// @Failure 500 internal error
// @router / [post]
func (u *UserController) Post() {
	var user *models.DemoUser
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		u.Logger.Infof("Parse RequestBody err: %v", err)
		u.ServeJSON(400, err.Error())
		return
	}

	uid, err := user.Add()
	if err != nil {
		u.Logger.Errorf("Add err: %v", err)
		u.ServeJSON(500, err.Error())
		return
	}

	u.ServeJSON(200, map[string]int{"uid": uid})
}

// Get get
// @Title Get
// @Description get user by uid
// @Param	uid		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.DemoUser
// @Failure 400 invalid uid
// @Failure 404 user not found
// @Failure 500 internal error
// @router /:uid [get]
func (u *UserController) Get() {
	uid, err := u.GetInt(":uid")
	if err != nil {
		u.Logger.Infof("GetInt err: %v", err)
		u.ServeJSON(400, err.Error())
		return
	}

	user := &models.DemoUser{ID: uid}
	err = user.Get()
	switch err {
	case models.ErrUNE:
		u.ServeJSON(404, err.Error())
		return
	case nil:
		u.ServeJSON(200, user)
		return
	default:
		u.Logger.Errorf("Get err: %v", err)
		u.ServeJSON(500, err.Error())
		return
	}
}

// Put update
// @Title Update
// @Description update the user
// @Param	uid		path 	int	true		"The uid you want to update"
// @Param	body		body 	models.DemoUser	true		"body for user content"
// @Success 200 {object} models.DemoUser
// @Failure 400 invalid uid/body
// @Failure 404 user not exist
// @Failure 500 internal error
// @router /:uid [put]
func (u *UserController) Put() {
	uid, err := u.GetInt(":uid")
	if err != nil {
		u.Logger.Infof("GetInt err: %v", err)
		u.ServeJSON(400, err.Error())
		return
	}

	user := &models.DemoUser{ID: uid}
	err = user.Get()
	switch err {
	case models.ErrUNE:
		u.ServeJSON(404, err.Error())
		return
	case nil:
		err = json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		if err != nil {
			u.Logger.Infof("Parse RequestBody err: %v", err)
			u.ServeJSON(400, err.Error())
			return
		}

		user.ID = uid
		err = user.Update()
		switch err {
		case nil:
			u.ServeJSON(200, user)
			return
		default:
			u.Logger.Errorf("Update err: %v", err)
			u.ServeJSON(500, err.Error())
			return
		}
	default:
		u.Logger.Errorf("Get err: %v", err)
		u.ServeJSON(500, err.Error())
		return
	}
}

// Delete delete
// @Title Delete
// @Description delete the user
// @Param	uid		path 	int	true		"The uid you want to delete"
// @Success 200 {json} {"uid": models.DemoUser.ID}
// @Failure 400 invalid uid
// @Failure 404 user not exist
// @Failure 500 internal error
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid, err := u.GetInt(":uid")
	user := &models.DemoUser{ID: uid}
	err = user.Delete()
	switch err {
	case models.ErrUNE:
		u.ServeJSON(404, err.Error())
		return
	case nil:
		u.ServeJSON(200, map[string]int{"uid": uid})
		return
	default:
		u.Logger.Errorf("Delete err: %v", err)
		u.ServeJSON(500, err.Error())
		return
	}
}
