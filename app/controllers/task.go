package controllers

import (
	"strconv"
	"viltasks/app/models"

	"github.com/revel/revel"
	"gopkg.in/robfig/cron.v2"
)

type Task struct {
	*revel.Controller
}

func (c Task) Index() revel.Result {
	return c.Render()
}

func (c Task) Delete() revel.Result {

	id := c.Params.Query.Get("id")
	i2, _ := strconv.Atoi(id)
	i3 := cron.EntryID(i2)
	models.Remove(i3)
	return c.Redirect(Task.List)
}

func (c Task) List() revel.Result {

	t := models.ListJob()
	return c.Render(t)
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

	//var t models.CronTask
	command := c.Params.Query.Get("command")
	name := c.Params.Query.Get("name")
	time := c.Params.Query.Get("time")
	desc := c.Params.Query.Get("desc")

	c.Validation.Required(name).Message("Nombre es requerido!")
	c.Validation.Required(desc).Message("Descripcion es requerido!")
	c.Validation.Required(command).Message("Comando es requerido!")
	c.Validation.Required(time).Message("Tiempo es requerido!")
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

	err := models.Addjob(t)
	if err != nil {

		c.Flash.Error(err.Error())
		return c.Redirect(Task.Create)
	}

	c.Flash.Success("Tarea realizada correctamente")
	return c.Redirect(Task.Create)
}

func (c Task) Create() revel.Result {
	return c.Render()
}

func (c Task) Clean() revel.Result {
	models.CleanFailedJobs()
	return c.Redirect(Task.Status)
}

func (c Task) Cleansuccess() revel.Result {
	models.CleanSuccessdJobs()
	return c.Redirect(Task.Status)
}
