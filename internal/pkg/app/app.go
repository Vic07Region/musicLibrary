package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"musicLibrary/internal/app/endpoint"
	"musicLibrary/internal/connector/songinfo"
	"musicLibrary/internal/database"
	"musicLibrary/internal/database/migrate"
	"musicLibrary/internal/lib/logger"
	"musicLibrary/internal/service"
	"os"
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
	a.l = logger.New()

	err := godotenv.Load(".env.local")
	if err != nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	//init config
	dbdriver := os.Getenv("DB_DRIVER")
	migrationsDIRS := os.Getenv("MIGRATION_DIRS")
	api_baseurl := os.Getenv("API_BASEURL")
	if dbdriver == "" {
		dbdriver = POSTGRES
	}
	cs_option := make(map[string]string)
	cs_option["host"] = os.Getenv("DB_HOST")
	cs_option["port"] = os.Getenv("DB_PORT")
	cs_option["sslmode"] = os.Getenv("DB_SSLMODE")
	cs_option["sslrootcert"] = os.Getenv("DB_ROOTSERT")
	a.l.Debug("cs_option", cs_option)
	cs := newCS(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		cs_option,
	)

	//init layers
	//database connection
	db, err := database.NewConnection(dbdriver, cs)
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
	song_info_service := songinfo.New(api_baseurl)
	//init service layer
	a.s = service.New(a.dbq, song_info_service, a.l)
	//init endpoint
	a.e = endpoint.New(a.s)

	a.gin = gin.Default()

	a.gin.GET("api/songs", a.e.FetchSongsHandler)
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
	log.Println(connectionString)
	return connectionString
}
