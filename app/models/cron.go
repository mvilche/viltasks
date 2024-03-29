package models

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/revel/revel"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Tz struct {
	Name string
}

type CronTask struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	gorm.Model
	Name                string `gorm:"size:255"`
	Description         string `gorm:"size:255"`
	Command             string `gorm:"size:255"`
	Time                string `gorm:"size:255"`
	CronId              string `gorm:"size:255"`
	Notification        bool
	Notification_failed bool
	Notification_email  string `gorm:"size:255"`
	Timezone            string `gorm:"size:255"`
}

func NewCronTask(t CronTask) *CronTask {
	return &CronTask{
		Name:                t.Name,
		Description:         t.Description,
		Command:             t.Command,
		Time:                t.Time,
		CronId:              t.CronId,
		Notification:        t.Notification,
		Notification_email:  t.Notification_email,
		Notification_failed: t.Notification_failed,
		Timezone:            t.Timezone,
	}
}

type FailedCronTask struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	gorm.Model
	Name   string `gorm:"size:255"`
	Output string `gorm:"size:255"`
	CronId string `gorm:"size:255"`
	Date   string `gorm:"size:255"`
}

func NewFailedCronTask(t FailedCronTask) *FailedCronTask {
	return &FailedCronTask{
		Name:   t.Name,
		CronId: t.CronId,
		Output: t.Output,
		Date:   t.Date,
	}
}

type SuccessCronTask struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	gorm.Model
	Name string `gorm:"size:255"`
	Date string `gorm:"size:255"`
}

func NewSuccessCronTask(t SuccessCronTask) *SuccessCronTask {
	return &SuccessCronTask{
		Name: t.Name,
		Date: t.Date,
	}
}

type CronTaskConfig struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	gorm.Model
	Name string `gorm:"size:255"`
}

func NewCronTaskConfig(t CronTaskConfig) *CronTaskConfig {
	return &CronTaskConfig{
		Name: t.Name,
	}
}

var instance *cron.Cron
var once sync.Once

func GetCron() *cron.Cron {

	once.Do(func() {
		instance = cron.New()
		revel.AppLog.Debug("Init Cron")
		db, _ := OpenSQL()
		db.Exec("delete from cron_task_configs")
	})
	return instance
}

func StartCron() error {
	c := GetCron()
	var gerr error
	var tc CronTaskConfig
	tc.Name = "CronStarted"
	tconfig := NewCronTaskConfig(tc)
	db, _ := OpenSQL()
	if err := db.Where("name = ?", tc.Name).First(&tconfig).Error; err == nil {
		gerr = errors.New("Cron ya fue iniciado")
		revel.AppLog.Error(gerr.Error())
		return gerr

	}

	var r []CronTask
	db.Find(&r)
	n := len(r)
	if n > 0 {
		db.Exec("delete from cron_tasks")

		for _, b := range r {
			Addjob(b)
		}

	}
	revel.AppLog.Infof("Iniciando sistema con " + strconv.FormatInt(int64(n), 10) + " tasks")
	c.Start()
	db.Exec("delete from cron_task_configs")
	if err := db.Create(&tconfig).Error; err != nil {
		revel.AppLog.Error("Error al insertar config")
	}

	return gerr
}

func StopCron() error {

	c := GetCron()
	var tc CronTaskConfig
	var gerr error
	tc.Name = "CronStopped"
	tconfig := NewCronTaskConfig(tc)
	db, _ := OpenSQL()

	if err := db.Where("name = ?", tc.Name).First(&tconfig).Error; err == nil {
		gerr = errors.New("Cron ya se encuentra detenido")
		revel.AppLog.Error(gerr.Error())
		return gerr
	}
	revel.AppLog.Debug("Stoping cron")
	if len(c.Entries()) > 0 {
		revel.AppLog.Debug("Se detectaron jobs activos, quitando antes de detener")
		for _, b := range c.Entries() {

			c.Remove(b.ID)
		}
		if len(c.Entries()) <= 0 {
			revel.AppLog.Debug("Sin jobs activos")
		}
	}
	c.Stop()
	db.Exec("delete from cron_task_configs")
	tc.Name = "CronStopped"
	if err := db.Create(&tconfig).Error; err != nil {
		revel.AppLog.Error("Error al insertar config")
	}
	return gerr
}

func Addjob(t CronTask) error {

	c := GetCron()

	var gerror error
	db, _ := OpenSQL()
	task := NewCronTask(t)
	if err := db.Where("name = ?", t.Name).First(&task).Error; err == nil {
		gerror = errors.New("Ya existe un cron con el nombre ingresado")
		revel.AppLog.Error(gerror.Error())
		return gerror
	}

	id, _ := c.AddFunc("CRON_TZ="+t.Timezone+" "+t.Time, func() {

		cmd := exec.Command("/bin/sh", "-c", t.Command)

		out, err := cmd.CombinedOutput()

		if err != nil {
			revel.AppLog.Error("Error ejecutando task")
			revel.AppLog.Debug("Task: " + t.Name)
			revel.AppLog.Error(string(out) + " - " + err.Error())
			var f FailedCronTask

			f.Name = t.Name
			f.Output = string(out) + " - " + err.Error()
			f.Date = time.Now().Format("2006-01-02 15:04:05")
			failed := NewFailedCronTask(f)

			if err := db.Create(&failed).Error; err != nil {
				revel.AppLog.Error(err.Error())

			}

			if err := db.Where("name = ? AND notification = ?", task.Name, true).First(&task).Error; err == nil {

				revel.AppLog.Debug("Se encontro activado notificaciones")
				err2 := SendNewEmail(&SuccessCronTask{}, task, failed, false)
				if err2 != nil {
					revel.AppLog.Error(err2.Error())
				}
			}

		} else {
			revel.AppLog.Debug("Task: " + t.Name)
			var sTask SuccessCronTask
			sTask.Name = t.Name
			sTask.Date = time.Now().Format("2006-01-02 15:04:05")
			succes := NewSuccessCronTask(sTask)

			if err := db.Where("name = ?", succes.Name).First(&succes).Error; err != nil {

				if err := db.Create(&succes).Error; err != nil {
					revel.AppLog.Error(err.Error())

				}
			} else {
				succes.Date = time.Now().Format("2006-01-02 15:04:05")
				if err := db.Save(&succes).Error; err != nil {
					revel.AppLog.Error(err.Error())

				}
			}

			if err := db.Where("name = ? AND notification = ? AND notification_failed = false", task.Name, true).First(&task).Error; err == nil {

				revel.AppLog.Debug("Se encontro activado notificaciones")

				err2 := SendNewEmail(succes, task, &FailedCronTask{}, true)
				if err2 != nil {
					revel.AppLog.Error(err2.Error())
				}
			}

			revel.AppLog.Debug(string(out))

		}
	})

	if !c.Entry(id).Valid() {
		gerror = errors.New("Cron ingresado no valido")
		revel.AppLog.Error("Cron ingresado no valido")
		c.Remove(id)
		return gerror
	} else {
		revel.AppLog.Debug("Cron agregado valido")
	}

	task.CronId = strconv.FormatInt(int64(id), 10)
	if err := db.Create(&task).Error; err != nil {
		revel.AppLog.Error(err.Error())

	} else {
		revel.AppLog.Debug("Cron cargado en database")
	}

	return gerror

}

func ListJob() []CronTask {

	var result []CronTask
	db, _ := OpenSQL()
	db.Find(&result)
	return result
}

func SuccesJob() []SuccessCronTask {

	var result []SuccessCronTask
	db, _ := OpenSQL()
	db.Find(&result)
	return result
}

func ListFailedJob() []FailedCronTask {

	var result []FailedCronTask
	db, _ := OpenSQL()
	db.Find(&result)
	return result
}

func ShowTZ() []Tz {

	var tlist []Tz
	files, err := ioutil.ReadDir("/usr/share/zoneinfo/")
	if err != nil {

	}

	for _, f := range files {
		files2, err2 := ioutil.ReadDir("/usr/share/zoneinfo/" + f.Name())
		if err2 == nil {
			for _, f2 := range files2 {
				var t Tz
				sname := strings.TrimSpace(f.Name() + "/" + f2.Name())
				t.Name = sname
				tlist = append(tlist, t)
			}

		}

	}

	return tlist
}

func Remove(id cron.EntryID) {
	c := GetCron()
	db, _ := OpenSQL()
	e := c.Entry(id)
	if e.Valid() {
		c.Remove(id)
		if err := db.Where("cron_id = ?", strconv.FormatInt(int64(id), 10)).Delete(CronTask{}).Error; err != nil {
			revel.AppLog.Info("Error al eliminar el cron", err.Error)
		} else {
			revel.AppLog.Info("Cron eliminado")
		}
	} else {

		revel.AppLog.Info("No se encontro cron para el id dado")
	}
}

func CleanFailedJobs() error {
	db, _ := OpenSQL()
	if err := db.Exec("delete from failed_cron_tasks").Error; err != nil {

		return err
	}

	return nil
}

func CleanSuccessdJobs() error {

	db, _ := OpenSQL()
	if err := db.Exec("delete from success_cron_tasks").Error; err != nil {

		return err
	}

	return nil
}
