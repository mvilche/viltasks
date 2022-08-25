package app

import (
	"os"
	"viltasks/app/models"

	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	revel.OnAppStart(ShowVersion)
	revel.OnAppStart(InitMigrations)
	//revel.OnAppStart(InitCron)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}

func ShowVersion() {

	revel.AppLog.Info("Autor: Martin Fabrizzio Vilche")
	revel.AppLog.Info("Version: 2.1.5")

}

func InitMigrations() {
	db, err := models.OpenSQL()
	if err == nil {

		revel.AppLog.Info("Start database migrations")

		//list models

		err := db.AutoMigrate(&models.CronTask{}, &models.CronTaskConfig{}, &models.FailedCronTask{}, &models.SuccessCronTask{}, &models.User{})

		if err != nil {

			revel.AppLog.Error("Database migrations fail\n" + err.Error())

		} else {
			//db.Model(&models.User{}).AddForeignKey("rol_id", "rols(id)", "RESTRICT", "RESTRICT")
			//db.Model(&models.Ticket{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
			//	db.Model(&models.Ticket{}).AddForeignKey("status_id", "statuses(id)", "RESTRICT", "RESTRICT")
			//db.Model(&models.Group{}).AddForeignKey("user", "users(group_refer)", "RESTRICT", "RESTRICT")
			//	db.Model(&models.Group{}).AddForeignKey("group_id", "users(group_id)", "CASCADE", "CASCADE")
			//// INSERT INTO "users" (name) VALUES ("non_existing");
			//// user -> User{Id: 112, Name: "non_existing"}
			revel.AppLog.Info("Dat6abase migration finish")
		}

	} else {

		revel.AppLog.Error("Erro open database ", err)
		os.Exit(1)
	}

	InitCron()
}

func InitCron() {

	err := models.StartCron()
	if err != nil {
		revel.AppLog.Debug(err.Error())
	}

}
