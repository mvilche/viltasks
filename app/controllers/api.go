package controllers

import (
	"viltasks/app/models"

	"github.com/revel/revel"
)

type Api struct {
	*revel.Controller
}

func (c Api) Index() revel.Result {

	c.Response.Status = 403
	return c.RenderJSON("Forbidden")
}

func (c Api) ListFailed() revel.Result {
	p := models.ListFailedJob()
	c.Response.Out.Header().Add("Access-Control-Allow-Origin", "*")
	return c.RenderJSON(p)
}

func (c Api) ListSuccess() revel.Result {
	p := models.SuccesJob()
	c.Response.Out.Header().Add("Access-Control-Allow-Origin", "*")
	return c.RenderJSON(p)
}

func (c Api) Check() revel.Result {

	if c.Session["user"] != "L53PoWkpXMS3c2IUtGjHGQ" {
		return c.Redirect(Auth.Index)
	}

	return nil
}

func init() {
	revel.InterceptMethod(Api.Check, revel.BEFORE)
	revel.InterceptMethod(Task.Check, revel.BEFORE)
}
