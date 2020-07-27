package controllers

import (
	"viltasks/app/models"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Alta() revel.Result {

	var task models.CronTask

	task.Name = "martin"
	task.Time = "*/1 * * * *"
	task.Command = "curl -v --fail http://googlegol.com.uy"
	models.Addjob(task)
	return c.Render()
}

func (c App) Stop() revel.Result {
	models.StopCron()
	return c.Render()
}

func (c App) Start() revel.Result {
	models.StartCron()
	return c.Render()
}
