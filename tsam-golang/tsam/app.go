package tsam

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	// Mysql Dialect for database connection
	_ "github.com/jinzhu/gorm/dialects/mysql"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/techlabs/swabhav/tsam/config"
	"github.com/techlabs/swabhav/tsam/event"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/security"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	env "github.com/techlabs/swabhav/tsam/config"
)

// RouterSpecifier Implement By Controller For Endpoint Register.
type RouterSpecifier interface {
	RegisterRoutes(router *mux.Router, exclude *[]*mux.Route)
}

// Controller is implemented by the controllers.
type Controller interface {
	RegisterRoutes(router *mux.Router)
}

// ModuleConfig needs to be implemented by every module.
type ModuleConfig interface {
	TableMigration(wg *sync.WaitGroup)
}

// App Struct For Start the tsam service.
type App struct {
	Name           string
	Router         *mux.Router
	DB             *gorm.DB
	Log            log.Logger
	Config         env.ConfReader
	Server         *http.Server
	WG             *sync.WaitGroup
	Auth           *security.Authentication
	EventPool      *event.Pool
	IsInProduction bool
}

// NewApp returns app.
func NewApp(name string, db *gorm.DB, log log.Logger, config env.ConfReader, wg *sync.WaitGroup,
	auth *security.Authentication, isProd bool) *App {
	return &App{
		Name:   name,
		DB:     db,
		Log:    log,
		Config: config,
		WG:     wg,
		Auth:   auth,
		// EventPool:      pool,
		IsInProduction: isProd,
	}
}

// InitializeRouter Register the route.
// # new router
// func (app *App) InitializeRouter(router []Handler) {
func (app *App) InitializeRouter(router []RouterSpecifier) {
	app.Log.Info(app.Name + " App Route initializing")

	// slice for storing routes which are to be excluded
	// comment below line # new router
	var excludeRoutes []*mux.Route
	// New Router instance Created and Header Set.
	app.Router = mux.NewRouter().StrictSlash(true)
	app.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	app.Router = app.Router.PathPrefix("/api/v1/tsam").Subrouter()
	// # new router
	// app.Router = app.Router.PathPrefix("/api/v1/tsam").Subrouter()

	app.RouterRegister(router, excludeRoutes)
	// # new router
	// app.HandlerRegister(router)
	app.initializeServer()
}

func (app *App) initializeServer() {
	headers := handlers.AllowedHeaders([]string{
		"Content-Type", "X-Total-Count", "token", "totalLifetimeValue",
	})
	methods := handlers.AllowedMethods([]string{
		http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete, http.MethodOptions,
	})
	origin := handlers.AllowedOrigins([]string{
		// "http://localhost:4200", "http://localhost:4201", "http://localhost:4202",
		// "https://admin.swabhavtechlabs.com", "http://admin.swabhavtechlabs.com",
		"*",
	})
	apiPort := app.getPort()
	app.Server = &http.Server{
		Addr:         ":" + apiPort,
		ReadTimeout:  time.Second * time.Duration(app.Config.GetInt64("HTTP_READ_TIMEOUT")),
		WriteTimeout: time.Second * time.Duration(app.Config.GetInt64("HTTP_WRITE_TIMEOUT")),
		IdleTimeout:  time.Second * time.Duration(app.Config.GetInt64("HTTP_IDLE_TIMEOUT")),
		Handler:      handlers.CORS(headers, methods, origin)(app.Router),
	}
	app.Log.Printf("Server Exposed On %s", apiPort)
}

// RouterRegister will register the specified routes.
// #niranjan: need to delete this after RegisterControllerRoutes is implemented throughout.
func (app *App) RouterRegister(routerSpecifiers []RouterSpecifier, excludeRoutes []*mux.Route) {
	// Router Registration
	for _, routerSpecifier := range routerSpecifiers {
		// Can't use go routines as gorilla mux doesn't support it.
		routerSpecifier.RegisterRoutes(app.Router.NewRoute().Subrouter(), &excludeRoutes)
	}
	app.Router.Use(app.Auth.RegisterRoutes(excludeRoutes))
}

// RegisterControllerRoutes will register the specified routes in controllers.
func (app *App) RegisterControllerRoutes(controllers []Controller) {
	// controllers registering routes.
	for _, controller := range controllers {
		controller.RegisterRoutes(app.Router.NewRoute().Subrouter())
	}
}

// MigrateTables Register Modules Tables.
func (app *App) MigrateTables(configs []ModuleConfig) {
	app.WG.Add(len(configs))
	for _, config := range configs {
		// Check this code.
		go func(config ModuleConfig) {
			config.TableMigration(app.WG)
			app.WG.Done()
		}(config)
	}
	app.WG.Wait()
	app.Log.Info("End of Migration")
}

// Start will start the app.
func (app *App) Start() error {
	if err := app.Server.ListenAndServe(); err != nil {
		app.Log.Error("=========listen and serve error========", err)
		return err
	}
	app.Log.Info("Server Running!")
	return nil
}

// Stop the App.
func (app *App) Stop() {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(10))
	defer cancel()
	// Stopping Server.
	err := app.Server.Shutdown(context)

	if err != nil {
		app.Log.Fatal("Fail to Stop Server...")
		return
	}
	app.Log.Info("Server Shutdown")
}

func (app *App) getPort() string {
	if app.IsInProduction {
		return os.Getenv("PORT")
	}
	return app.Config.GetString(config.PORT)
}

func getConnectionString(conf config.ConfReader) (url string) {
	url += conf.GetString(config.DBUser)
	url += ":"
	url += conf.GetString(config.DBPass)
	url += "@tcp("
	url += conf.GetString(config.DBHost)
	url += ":"
	url += conf.GetString(config.DBPort)
	url += ")/"
	url += conf.GetString(config.DBName)
	url += "?charset=utf8&parseTime=true"
	return
}

// NewDBConnection Return DB Instace
func NewDBConnection(log log.Logger, conf config.ConfReader) *gorm.DB {
	url := getConnectionString(conf)
	log.Info("HERE IS THE OPEN URL :", url)
	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Error(err.Error())
		log.Fatal("Failed to create database instance.")
		return nil
	}
	db.LogMode(true)
	// db.SetLogger(log)
	db.BlockGlobalUpdate(true)
	return db
}
