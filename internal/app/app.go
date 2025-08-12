package app

import (
	"log"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/grvbrk/async0_server/internal/auth"
	adminHandler "github.com/grvbrk/async0_server/internal/handlers/admin"
	"github.com/grvbrk/async0_server/internal/middlewares"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/store/admin"
	"github.com/grvbrk/async0_server/migrations"
)

var (
	authKey            = securecookie.GenerateRandomKey(64)
	encryptionKey      = securecookie.GenerateRandomKey(32)
	adminAuthKey       = securecookie.GenerateRandomKey(64)
	adminEncryptionKey = securecookie.GenerateRandomKey(32)
)

type Application struct {
	Logger              *log.Logger
	Oauth               *auth.GoogleOauth
	AdminOauth          *auth.AdminGoogleOauth
	MiddlewareHandler   *middlewares.MiddlwareHandler
	AdminProblemHandler *adminHandler.AdminProblemHandler
	AdminListHandler    *adminHandler.AdminListHandler
	AdminTopicHandler   *adminHandler.AdminTopicHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "LOGGING: ", log.Ldate|log.Ltime)
	adminLogger := log.New(os.Stdout, "ADMIN LOGGING: ", log.Ldate|log.Ltime)
	sessionStore := sessions.NewCookieStore(authKey, encryptionKey)
	adminSessionStore := sessions.NewCookieStore(adminAuthKey, adminEncryptionKey)
	middlewareHandler := middlewares.NewMiddlewareHandler(logger, sessionStore)

	pgDB, err := store.ConnectPGDB()
	if err != nil {
		logger.Println("Error connecting to db")
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, "db")
	if err != nil {
		logger.Println("PANIC: Postgresql migration failed, exiting...")
		panic(err)
	}

	logger.Println("Database migrated...")

	// stores
	userStore := store.NewPostgresUserStore(pgDB)

	// admin stores
	// adminUserStore := admin.NewPostgresAdminUserStore(pgDB)
	adminProblemStore := admin.NewPostgresAdminProblemStore(pgDB)
	adminListStore := admin.NewPostgresAdminListStore(pgDB)
	adminTopicStore := admin.NewPostgresAdminTopicStore(pgDB)

	oauth, err := auth.NewGoogleOauth(logger, sessionStore, userStore)
	if err != nil {
		return nil, err
	}

	adminOauth, err := auth.NewAdminGoogleOauth(adminLogger, adminSessionStore, userStore)
	if err != nil {
		return nil, err
	}

	// handlers

	// admin handlers
	adminProblemHandler := adminHandler.NewAdminProblemHandler(adminProblemStore, adminLogger, adminOauth)
	adminListHandler := adminHandler.NewAdminListHandler(adminListStore, adminLogger, adminOauth)
	adminTopicHandler := adminHandler.NewAdminTopicHandler(adminTopicStore, adminLogger, adminOauth)

	app := &Application{
		Logger:              logger,
		Oauth:               oauth,
		AdminOauth:          adminOauth,
		MiddlewareHandler:   middlewareHandler,
		AdminProblemHandler: adminProblemHandler,
		AdminListHandler:    adminListHandler,
		AdminTopicHandler:   adminTopicHandler,
	}

	return app, nil
}
