package models

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/Kamva/mgm"
	"github.com/revel/revel"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/robfig/cron.v2"
)

type CronTask struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Description      string `json:"description" bson:"description"`
	Command          string `json:"command" bson:"command"`
	Time             string `json:"time" bson:"time"`
	CronId           string `json:"cronid" bson:"cronid"`
}

func NewCronTask(t CronTask) *CronTask {
	return &CronTask{
		Name:        t.Name,
		Description: t.Description,
		Command:     t.Command,
		Time:        t.Time,
		CronId:      t.CronId,
	}
}

type FailedCronTask struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Output           string `json:"output" bson:"output"`
	CronId           string `json:"cronid" bson:"cronid"`
	Date             string `json:"date" bson:"date"`
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
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Date             string `json:"date" bson:"date"`
}

func NewSuccessCronTask(t SuccessCronTask) *SuccessCronTask {
	return &SuccessCronTask{
		Name: t.Name,
		Date: t.Date,
	}
}

type CronTaskConfig struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
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
		mgm.Coll(&CronTaskConfig{}).Drop(mgm.Ctx())
	})
	return instance
}

func StartCron() error {
	c := GetCron()
	var gerr error
	var tc CronTaskConfig
	tc.Name = "CronStarted"
	tconfig := NewCronTaskConfig(tc)

	err := mgm.Coll(tconfig).First(bson.M{"name": tc.Name}, tconfig)
	if err == nil {
		gerr = errors.New("Cron ya fue iniciado")
		revel.AppLog.Error(gerr.Error())
		return gerr
	}

	result := []CronTask{}
	mgm.Coll(&CronTask{}).SimpleFind(&result, bson.D{{}})
	n := len(result)
	if n > 0 {
		mgm.Coll(&CronTask{}).Drop(mgm.Ctx())
		for _, b := range result {
			Addjob(b)
		}

	}
	revel.AppLog.Infof("Iniciando sistema con " + strconv.FormatInt(int64(n), 10) + " tasks")
	c.Start()
	mgm.Coll(&CronTaskConfig{}).Drop(mgm.Ctx())
	errr := mgm.Coll(tconfig).Create(tconfig)
	if errr != nil {
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

	err := mgm.Coll(tconfig).First(bson.M{"name": tc.Name}, tconfig)
	if err == nil {
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
	mgm.Coll(&CronTaskConfig{}).Drop(mgm.Ctx())
	tc.Name = "CronStopped"
	mgm.Coll(tconfig).Create(&tc)
	return gerr
}

func Addjob(t CronTask) error {

	c := GetCron()
	var gerror error

	task := NewCronTask(t)
	err := mgm.Coll(task).First(bson.M{"name": t.Name}, task)
	if err == nil {
		gerror = errors.New("Ya existe un cron con el nombre ingresado")
		revel.AppLog.Error(gerror.Error())
		return gerror
	}

	id, _ := c.AddFunc(t.Time, func() {

		out, err := exec.Command(t.Command).Output()

		if err != nil {
			revel.AppLog.Error("Error ejecutando task")
			revel.AppLog.Debug("Task: " + t.Name)
			revel.AppLog.Error(err.Error())
			var f FailedCronTask

			f.Name = t.Name
			f.Output = string(out) + " - " + err.Error()
			f.Date = time.Now().Format("2006-01-02 15:04:05")
			failed := NewFailedCronTask(f)

			err := mgm.Coll(failed).Create(failed)
			if err != nil {
				revel.AppLog.Error(err.Error())

			}
		} else {
			revel.AppLog.Debug("Task: " + t.Name)
			var sTask SuccessCronTask
			sTask.Name = t.Name
			sTask.Date = time.Now().Format("2006-01-02 15:04:05")
			succes := NewSuccessCronTask(sTask)

			find := mgm.Coll(succes).First(bson.M{"name": succes.Name}, succes)
			if find != nil {

				err2 := mgm.Coll(succes).Create(succes)
				if err2 != nil {
					revel.AppLog.Error(err.Error())
				}

			} else {
				succes.Date = time.Now().Format("2006-01-02 15:04:05")
				err := mgm.Coll(succes).Update(succes)
				if err != nil {
					revel.AppLog.Error(err.Error())
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
		revel.AppLog.Info("Cron agregado valido")
	}

	task.CronId = strconv.FormatInt(int64(id), 10)
	err2 := mgm.Coll(task).Create(task)
	if err2 != nil {
		revel.AppLog.Error(err.Error())
	} else {
		revel.AppLog.Debug("Cron cargado en mongo")
	}

	return gerror

}

func ListJob() []CronTask {

	result := []CronTask{}
	mgm.Coll(&CronTask{}).SimpleFind(&result, bson.D{{}})
	return result
}

func SuccesJob() []SuccessCronTask {

	result := []SuccessCronTask{}
	mgm.Coll(&SuccessCronTask{}).SimpleFind(&result, bson.D{{}})
	return result
}

func ListFailedJob() []FailedCronTask {

	result := []FailedCronTask{}
	mgm.Coll(&FailedCronTask{}).SimpleFind(&result, bson.D{{}})
	return result
}

func Remove(id cron.EntryID) {
	c := GetCron()
	e := c.Entry(id)
	if e.Valid() {
		c.Remove(id)
		t := &CronTask{}
		mgm.Coll(t).DeleteOne(mgm.Ctx(), bson.M{"cronid": strconv.FormatInt(int64(id), 10)})
		revel.AppLog.Info("Cron eliminado")
	} else {

		revel.AppLog.Info("No se encontro cron para el id dado")
	}
}

func CleanFailedJobs() error {

	return mgm.Coll(&FailedCronTask{}).Drop(mgm.Ctx())
}

func CleanSuccessdJobs() error {

	return mgm.Coll(&SuccessCronTask{}).Drop(mgm.Ctx())
}

func Entry() {

	c := GetCron()
	fmt.Println(c.Entries())
}
