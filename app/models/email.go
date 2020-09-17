package models

import (
	"crypto/tls"
	"os"
	"strconv"

	"github.com/revel/config"
	"github.com/revel/revel"
	"gopkg.in/gomail.v2"
)

// SSL/TLS Email Example

func SendNewEmail(suc *SuccessCronTask, cron *CronTask, fcron *FailedCronTask, ok bool) error {

	c, err := config.ReadDefault("./conf/app.conf")
	if err != nil {
		revel.AppLog.Error(err.Error())
		os.Exit(1)
	}

	var host, _ = c.String("email", "mail.host")
	var sport, _ = c.String("email", "mail.port")
	port, _ := strconv.ParseInt(sport, 10, 64)
	var user, _ = c.String("email", "mail.user")
	var password, _ = c.String("email", "mail.password")
	var t, _ = c.String("email", "mail.disable.tls")
	disabletls, _ := strconv.ParseBool(t)
	var s string
	var subject string

	if ok {

		s, _, subject, _ = RenderHtml(cron, suc, fcron, true)

	} else {

		_, s, _, subject = RenderHtml(cron, suc, fcron, false)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", user)
	revel.AppLog.Debug("from: " + user)
	m.SetHeader("To", cron.Notification_email)
	revel.AppLog.Debug("to: " + cron.Notification_email)
	m.SetHeader("Subject", subject)
	revel.AppLog.Debug("subject: " + subject)
	m.SetBody("text/html", s)
	revel.AppLog.Debug("body: " + s)
	// Attach some file
	//m.Attach("myfile1.pdf")

	d := gomail.NewDialer(host, int(port), user, password)
	d.TLSConfig = &tls.Config{ServerName: host, InsecureSkipVerify: disabletls}
	revel.AppLog.Debug("tls :" + t)

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		revel.AppLog.Error(err.Error())
		return err
	}
	revel.AppLog.Debug("send ok!")
	return nil
}

func RenderHtml(cron *CronTask, success *SuccessCronTask, fcron *FailedCronTask, ok bool) (string, string, string, string) {

	sOk := `<h3>Tarea ejecutada correctamente</h3>
	<h2>Tarea: </h2><p>` + cron.Name + `</p>
	<h4>Descripción: </h4><p>` + cron.Description + `</p>
	<h4>Hora de ejecución: </h4><p>` + success.Date + `</p>
	</br>---------------------------------------------------</br>
	<p>viltasks</p>`

	sError := `<h3>Tarea fallida!!!</h3>
	<h2>Tarea: </h2><p>` + fcron.Name + `</p>
	<h4>Descripción: </h4><p>` + cron.Description + `</p>
	<h4>Hora de ejecución: </h4><p>` + fcron.Date + `</p>
	<h4>Problema: </h4><p>` + fcron.Output + `</p>
	</br>---------------------------------------------------</br>
	<p>viltasks</p>`

	subjectOk := "Notificación de ejecución exitosa - Tarea: " + cron.Name + ""

	subjectError := "Notificación de ejecución fallida!! - Tarea: " + cron.Name + ""

	return sOk, sError, subjectOk, subjectError

}
