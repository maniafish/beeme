package controllers

import (
	"beeme/models"
	"encoding/json"

	"github.com/astaxie/beego"
)

// ObjectController Operations about object
type ObjectController struct {
	beego.Controller
}

// Post post
// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (o *ObjectController) Post() {
	var ob models.Object
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	objectid := models.AddOne(ob)
	o.Data["json"] = map[string]string{"ObjectId": objectid}
	o.ServeJSON()
}

// Get get
// @Title Get
// @Description find object by objectid
// @Param	objectid		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectid is empty
// @router /:objectid [get]
func (o *ObjectController) Get() {
	objectID := o.Ctx.Input.Param(":objectid")
	if objectID != "" {
		ob, err := models.GetOne(objectID)
		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			o.Data["json"] = ob
		}
	}
	o.ServeJSON()
}

// GetAll getall
// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *ObjectController) GetAll() {
	obs := models.GetAll()
	o.Data["json"] = obs
	o.ServeJSON()
}

// Put update
// @Title Update
// @Description update the object
// @Param	objectid		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectid is empty
// @router /:objectid [put]
func (o *ObjectController) Put() {
	objectID := o.Ctx.Input.Param(":objectid")
	var ob models.Object
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)

	err := models.Update(objectID, ob.Score)
	if err != nil {
		o.Data["json"] = err.Error()
	} else {
		o.Data["json"] = "update success!"
	}
	o.ServeJSON()
}

// Delete delete
// @Title Delete
// @Description delete the object
// @Param	objectid		path 	string	true		"The objectid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectid is empty
// @router /:objectid [delete]
func (o *ObjectController) Delete() {
	objectID := o.Ctx.Input.Param(":objectid")
	models.Delete(objectID)
	o.Data["json"] = "delete success!"
	o.ServeJSON()
}
