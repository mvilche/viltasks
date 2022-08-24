package controllers

import (
	"strconv"
	"viltasks/app/models"

	"github.com/revel/revel"
	"github.com/robfig/cron/v3"
)

type Task struct {
	*revel.Controller
}

func (c Task) Index() revel.Result {
	return c.Render()
}

func (c Task) Delete() revel.Result {

	id := c.Params.Get("id")
	i2, _ := strconv.Atoi(id)
	i3 := cron.EntryID(i2)
	models.Remove(i3)
	return c.Redirect(Task.List)
}

func (c Task) List() revel.Result {

	t := models.ListJob()
	return c.Render(t)
}

func (c Task) Sintaxis() revel.Result {
	return c.Render()
}

func (c Task) Status() revel.Result {

	t := models.ListFailedJob()
	s := models.SuccesJob()

	if len(t) <= 0 {

		return c.Render(s)
	}

	return c.Render(t)
}

func (c Task) CreateTask() revel.Result {

	command := c.Params.Form.Get("command")
	name := c.Params.Form.Get("name")
	time := c.Params.Form.Get("time")
	desc := c.Params.Form.Get("desc")
	email := c.Params.Form.Get("email")
	noti := c.Params.Form.Get("notificacion")
	notiF := c.Params.Form.Get("notificacion_failed")
	tz := c.Params.Form.Get("timezone")

	c.Validation.Required(name).Message("Nombre es requerido!")
	c.Validation.Required(desc).Message("Descripcion es requerido!")
	c.Validation.Required(command).Message("Comando es requerido!")
	c.Validation.Required(time).Message("Tiempo es requerido!")

	if revel.ToBool(noti) {

		c.Validation.Required(email).Message("Email es requerido!")
		c.Validation.Email(email).Message("Email no es valido!")

	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Task.Create)
	}

	var t models.CronTask
	t.Command = command
	t.Name = name
	t.Description = desc
	t.Time = time
	t.Notification = revel.ToBool(noti)
	t.Notification_failed = revel.ToBool(notiF)
	t.Notification_email = email
	t.Timezone = tz

	err := models.Addjob(t)
	if err != nil {

		c.Flash.Error(err.Error())
		return c.Redirect(Task.Create)
	}

	c.Flash.Success("Tarea realizada correctamente")
	return c.Redirect(Task.Create)
}

func (c Task) Create() revel.Result {
	tz := models.ShowTZ()
	if len(tz) < 0 {

		revel.AppLog.Error("Error al obtener timezone del sistema operativo")
	}

	return c.Render(tz)
}

func (c Task) Clean() revel.Result {
	models.CleanFailedJobs()
	return c.Redirect(Task.Status)
}

func (c Task) Cleansuccess() revel.Result {
	models.CleanSuccessdJobs()
	return c.Redirect(Task.Status)
}

func (c Task) Check() revel.Result {

	if c.Session["user"] != "L53PoWkpXMS3c2IUtGjHGQ" {
		return c.Redirect(Auth.Index)
	}

	return nil
}
