package app

import (
	"database/sql"
	"fmt" //nolint:gci
	"os"
	"strconv"
	"time"

	"github.com/Vic07Region/musicLibrary/docs"
	"github.com/Vic07Region/musicLibrary/internal/app/endpoint"
	"github.com/Vic07Region/musicLibrary/internal/connector/songinfo"
	"github.com/Vic07Region/musicLibrary/internal/database"
	"github.com/Vic07Region/musicLibrary/internal/database/migrate"
	"github.com/Vic07Region/musicLibrary/internal/lib/csmaker"
	"github.com/Vic07Region/musicLibrary/internal/lib/logger"
	"github.com/Vic07Region/musicLibrary/internal/service"
	"github.com/gin-gonic/gin" //nolint:gci
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	POSTGRES = "postgres"
)

type App struct {
	db  *sql.DB
	dbq database.Storage
	s   service.MusicService
	e   *endpoint.Endpoint
	l   *logger.Logger
	gin *gin.Engine
}

func New() (*App, error) {
	a := &App{}
	//init logger
	a.l = logger.New()

	//init dotenv
	err := godotenv.Load(".env.local")
	if err != nil {
		err = godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	//init config
	production := os.Getenv("PRODUCTION")
	var prod bool
	prod, err = strconv.ParseBool(production)
	if err != nil {
		return nil, fmt.Errorf("PRODUCTION .env param wrong (BOOL)")
	}

	debugParam := os.Getenv("DEBUG")
	var debug bool
	debug, err = strconv.ParseBool(debugParam)
	if err != nil {
		return nil, fmt.Errorf("DEBUG .env param wrong (BOOL)")
	}

	dbdriver := os.Getenv("DB_DRIVER")
	migrationsDIRS := os.Getenv("MIGRATION_DIRS")
	apiBaseurl := os.Getenv("API_BASEURL")

	maxConnEnv := os.Getenv("DB_MAX_CONN")
	macIdleEnv := os.Getenv("DB_MAX_IDLE")
	maxLifetimeEnv := os.Getenv("DB_MAX_LIFETIME")

	var maxCon, maxIdle int
	var maxLifetime time.Duration

	if maxConnEnv != "" {
		maxCon, err = strconv.Atoi(maxConnEnv)
		if err != nil {
			return nil, fmt.Errorf("DB_MAX_CONN param wrong (INT)")
		}
	} else {
		maxCon = 0
	}

	if macIdleEnv != "" {
		maxIdle, err = strconv.Atoi(macIdleEnv)
		if err != nil {
			return nil, fmt.Errorf("DB_MAX_IDLE param wrong (INT)")
		}
	} else {
		maxIdle = 5
	}

	if maxLifetimeEnv != "" {
		tm, err := strconv.Atoi(macIdleEnv)
		if err != nil {
			return nil, fmt.Errorf("DB_MAX_LIFETIME param wrong (INT)")
		}
		maxLifetime = time.Duration(tm) * time.Minute
	} else {
		maxLifetime = 0
	}

	if dbdriver == "" {
		dbdriver = POSTGRES
	}
	//make connection string option
	csOption := make(map[string]string)
	csOption["host"] = os.Getenv("DB_HOST")
	csOption["port"] = os.Getenv("DB_PORT")
	csOption["sslmode"] = os.Getenv("DB_SSLMODE")
	csOption["sslrootcert"] = os.Getenv("DB_ROOTSERT")
	//make connection string
	cs := csmaker.MakeConnectionString(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		csOption,
	)

	//init layers
	//database connection
	a.db, err = database.NewConnection(database.ConnectionParams{
		DbDriver:         dbdriver,
		ConnectionString: cs,
		MaxOpenConns:     maxCon,
		MsxIdleConns:     maxIdle,
		MaxLifeTime:      maxLifetime,
	})
	if err != nil {
		return nil, err
	}

	//migrator
	err = migrate.ApplyMigrations(a.db, migrationsDIRS, dbdriver)
	if err != nil {
		return nil, err
	}

	//init db queries storage
	a.dbq = database.NewStorage(a.db, a.l, debug)
	//init third api service
	songInfoService := songinfo.New(apiBaseurl, a.l)
	//init service layer
	a.s = service.New(a.dbq, songInfoService, a.l, debug)
	//init endpoint
	a.e = endpoint.New(a.s, a.l)

	//set gin mode
	if prod {
		gin.SetMode(gin.ReleaseMode)
	}
	//init gin engine
	a.gin = gin.Default()

	//set swagger basePath
	docs.SwaggerInfo.BasePath = "/api/v1"

	//register handlers
	eg := a.gin.Group("/api/v1")
	{
		eg.GET("/songs", a.e.FetchSongsHandler)
		eg.GET("/songs/:id", a.e.FetchSongTextHandler)
		eg.DELETE("/songs/:id", a.e.DeleteSongHandler)
		eg.PATCH("/songs/:id", a.e.UpdateSongHandler)
		eg.PATCH("/songs/:id/verse", a.e.UpdateSongVerseHandler)
		eg.POST("/songs/new", a.e.NewSongHandler)
	}
	//third route
	a.gin.GET("/info", a.e.TestHandler)
	a.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return a, nil
}

func (a *App) Run() error {
	defer a.db.Close()

	//init host
	hostOption := os.Getenv("APP_HOST")
	if hostOption == "" {
		hostOption = ":8080"
	}

	a.l.Info("Start gin server", "host", hostOption)
	err := a.gin.Run(hostOption)
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	return nil
}
