package controllers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"viltasks/app/models"

	"github.com/revel/revel"
)

type Auth struct {
	*revel.Controller
}

func getCredentials(data string) (username, password string, err error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", "", err
	}
	strData := strings.Split(string(decodedData), ":")
	username = strData[0]
	password = strData[1]
	return
}

func (c Auth) Index() revel.Result {
	// The auth data is sent in the headers and will have a value of "Basic XXX" where XXX is base64 encoded data
	if auth := c.Request.Header.Get("Authorization"); auth != "" {
		// Split up the string to get just the data, then get the credentials
		username, password, err := getCredentials(strings.Split(auth, " ")[1])
		if err != nil {
			return c.RenderError(err)
		}

		u, err := models.GetUser()
		if err != nil {
			return c.RenderError(err)
		}
		if username != u.Username || password != u.Password {
			c.Response.Status = http.StatusUnauthorized
			c.Response.Out.Header().Set("WWW-Authenticate", `Basic realm="revel"`)
			return c.RenderError(errors.New("401: Not authorized"))
		}
		c.Session["user"] = "L53PoWkpXMS3c2IUtGjHGQ"
		return c.Redirect(App.Index)
	} else {
		c.Response.Status = http.StatusUnauthorized
		c.Response.Out.Header().Set("WWW-Authenticate", `Basic realm="revel"`)
		return c.RenderError(errors.New("401: Not authorized"))
	}
}

func (c Auth) Logout() revel.Result {
	revel.NewSessionCookieEngine()
	return c.Redirect(Auth.Index)
}
