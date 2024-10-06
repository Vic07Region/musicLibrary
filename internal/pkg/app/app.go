package app

import (
	"fmt"
	docs "github.com/Vic07Region/musicLibrary/docs"
	"github.com/Vic07Region/musicLibrary/internal/app/endpoint"
	"github.com/Vic07Region/musicLibrary/internal/connector/songinfo"
	"github.com/Vic07Region/musicLibrary/internal/database"
	"github.com/Vic07Region/musicLibrary/internal/database/migrate"
	"github.com/Vic07Region/musicLibrary/internal/lib/logger"
	"github.com/Vic07Region/musicLibrary/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	POSTGRES = "postgres"
)

type App struct {
	dbq *database.Queries
	e   *endpoint.Endpoint
	s   *service.Service
	l   *logger.Logger
	gin *gin.Engine
}

func New() (*App, error) {
	a := &App{}

	err := godotenv.Load(".env.local")
	if err != nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	//init config
	debug_param := os.Getenv("DEBUG")
	var debug bool
	if debug_param == "TRUE" || debug_param == "true" || debug_param == "True" {
		debug = true
	}
	a.l = logger.New(debug)
	a.l.Info("DEBUG LEVEL", debug)
	dbdriver := os.Getenv("DB_DRIVER")
	migrationsDIRS := os.Getenv("MIGRATION_DIRS")
	api_baseurl := os.Getenv("API_BASEURL")

	max_conn_env := os.Getenv("DB_MAX_CONN")
	mac_idle_env := os.Getenv("DB_MAX_IDLE")
	max_lifetime_env := os.Getenv("DB_MAX_LIFETIME")

	var max_con, max_idle int
	var max_lifetime time.Duration
	if max_conn_env != "" {
		max_con, err = strconv.Atoi(max_conn_env)
		if err != nil {
			return nil, fmt.Errorf("Max conn param wrong")
		}
	} else {
		max_con = 25
	}

	if mac_idle_env != "" {
		max_idle, err = strconv.Atoi(mac_idle_env)
		if err != nil {
			return nil, fmt.Errorf("Max idle param wrong")
		}
	} else {
		max_idle = 5
	}

	if max_lifetime_env != "" {
		tm, err := strconv.Atoi(mac_idle_env)
		if err != nil {
			return nil, fmt.Errorf("Max lifetime param wrong")
		}
		max_lifetime = time.Duration(tm) * time.Minute
	} else {
		max_lifetime = 5 * time.Minute
	}

	if dbdriver == "" {
		dbdriver = POSTGRES
	}
	cs_option := make(map[string]string)
	cs_option["host"] = os.Getenv("DB_HOST")
	cs_option["port"] = os.Getenv("DB_PORT")
	cs_option["sslmode"] = os.Getenv("DB_SSLMODE")
	cs_option["sslrootcert"] = os.Getenv("DB_ROOTSERT")
	cs := newCS(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		cs_option,
	)

	//init layers
	//database connection
	db, err := database.NewConnection(database.ConnectionParams{
		DbDriver:         dbdriver,
		ConnectionString: cs,
		MaxOpenConns:     max_con,
		MsxIdleConns:     max_idle,
		MaxLifeTime:      max_lifetime,
	})
	if err != nil {
		return nil, err
	}

	//migrator
	err = migrate.ApplyMigrations(db, migrationsDIRS, dbdriver)
	if err != nil {
		return nil, err
	}
	//init db queries storage
	a.dbq = database.NewStorage(db, a.l)
	//init third api service
	song_info_service := songinfo.New(api_baseurl, a.l)
	//init service layer
	a.s = service.New(a.dbq, song_info_service, a.l)
	//init endpoint
	a.e = endpoint.New(a.s, a.l)

	a.gin = gin.Default()
	//swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	eg := a.gin.Group("/api/v1")
	{
		eg.GET("/songs", a.e.FetchSongsHandler)
		eg.GET("/songs/:id", a.e.FetchSongTextHandler)
		eg.DELETE("/songs/:id", a.e.DeleteSongHandler)
		eg.PATCH("/songs/:id", a.e.UpdateSongHandler)
		eg.POST("/songs/new", a.e.NewSongHandler)
	}
	a.gin.GET("/info", a.e.TestHandler)
	a.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return a, nil
}

func (a *App) Run() error {

	host_option := os.Getenv("APP_HOST")
	if host_option == "" {
		host_option = ":8080"
	}
	a.l.Info("Server running on", "addr", host_option)
	err := a.gin.Run(host_option)
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}

func newCS(user, password, dbName string, options ...map[string]string) string {
	host := "localhost"
	port := "5432"
	sslmode := "disable"
	sslrootcert := ""

	if len(options) > 0 {
		for _, option := range options {
			if val, ok := option["host"]; ok && val != "" {
				host = val
			}
			if val, ok := option["port"]; ok && val != "" {
				port = val
			}
			if val, ok := option["sslmode"]; ok && val != "" {
				sslmode = val
			}
			if val, ok := option["sslrootcert"]; ok && val != "" {
				sslrootcert = fmt.Sprintf("sslrootcert=%s", val)
			}
		}
	}

	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s %s",
		user, password, host, port, dbName, sslmode, sslrootcert)
	return connectionString
}
