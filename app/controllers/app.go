package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {

	return c.Render()
}

func (c App) Check() revel.Result {

	if c.Session["user"] != "L53PoWkpXMS3c2IUtGjHGQ" {
		return c.Redirect(Auth.Index)
	}

	return nil
}

func init() {
	revel.InterceptMethod(App.Check, revel.BEFORE)
	revel.InterceptMethod(Task.Check, revel.BEFORE)
}

/*

func (c App) Stop() revel.Result {
	models.StopCron()
	return c.Render()
}

func (c App) Start() revel.Result {
	models.StartCron()
	return c.Render()
}
*/
