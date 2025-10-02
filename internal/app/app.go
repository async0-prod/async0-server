package app

import (
	"context"
	"log"
	"os"

	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/handlers"
	adminHandler "github.com/grvbrk/async0_server/internal/handlers/admin"
	"github.com/grvbrk/async0_server/internal/middlewares"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/store/admin"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

// var (
// 	authKey            = securecookie.GenerateRandomKey(64)
// 	encryptionKey      = securecookie.GenerateRandomKey(32)
// 	adminAuthKey       = securecookie.GenerateRandomKey(64)
// 	adminEncryptionKey = securecookie.GenerateRandomKey(32)
// )

type Application struct {
	Logger      *log.Logger
	redisClient *redis.Client

	Oauth      *auth.GoogleOauth
	AdminOauth *auth.AdminGoogleOauth

	MiddlewareHandler     *middlewares.MiddlewareHandler
	UserProblemHandler    *handlers.ProblemHandler
	UserListHandler       *handlers.ListHandler
	UserTestcaseHandler   *handlers.TestcaseHandler
	UserSubmissionHandler *handlers.SubmissionHandler
	UserTopicHandler      *handlers.TopicHandler

	AdminProblemHandler *adminHandler.AdminProblemHandler
	AdminListHandler    *adminHandler.AdminListHandler
	AdminTopicHandler   *adminHandler.AdminTopicHandler

	UserAnalyticsHandler *handlers.AnalyticsHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "LOGGING: ", log.Ldate|log.Ltime)
	adminLogger := log.New(os.Stdout, "ADMIN LOGGING: ", log.Ldate|log.Ltime)

	pgDB, err := store.ConnectPGDB()
	if err != nil {
		logger.Println("Error connecting to db")
		return nil, err
	}

	// err = store.MigrateFS(pgDB, migrations.FS, "db")
	// if err != nil {
	// 	logger.Println("PANIC: Postgresql migration failed, exiting...")
	// 	panic(err)
	// }

	// logger.Println("Database migrated...")

	redisClient, err := store.ConnectRedis()
	if err != nil {
		logger.Println("PANIC: Redis connection failed, exiting...")
		panic(err)
	}

	logger.Println("Redis connected...")

	sessionStore, err := redisstore.NewRedisStore(context.Background(), redisClient)
	if err != nil {
		logger.Println("PANIC: Redis session store failed, exiting...")
		panic(err)
	}

	logger.Println("Redis session store connected...")

	middlewareHandler := middlewares.NewMiddlewareHandler(logger, sessionStore)

	// user stores
	userStore := store.NewPostgresUserStore(pgDB)
	problemStore := store.NewPostgresProblemStore(pgDB)
	listStore := store.NewPostgresListStore(pgDB)
	testcaseStore := store.NewPostgresTestcaseStore(pgDB)
	submissionStore := store.NewPostgresSubmissionStore(pgDB)
	topicStore := store.NewPostgresTopicStore(pgDB)

	// admin stores
	adminProblemStore := admin.NewPostgresAdminProblemStore(pgDB)
	adminListStore := admin.NewPostgresAdminListStore(pgDB)
	adminTopicStore := admin.NewPostgresAdminTopicStore(pgDB)

	// analytics store
	analyticsStore := store.NewPostgresAnalyticsStore(pgDB)

	oauth, err := auth.NewGoogleOauth(logger, sessionStore, userStore)
	if err != nil {
		return nil, err
	}

	adminOauth, err := auth.NewAdminGoogleOauth(adminLogger, sessionStore, userStore)
	if err != nil {
		return nil, err
	}

	// user handlers
	userProblemHandler := handlers.NewProblemHandler(problemStore, logger, oauth)
	userListHandler := handlers.NewListHandler(listStore, logger, oauth)
	userTestcaseHandler := handlers.NewTestcaseHandler(testcaseStore, logger, oauth)
	userSubmissionHandler := handlers.NewSubmissionHandler(submissionStore, testcaseStore, logger, oauth)
	userTopicHandler := handlers.NewTopicHandler(topicStore, logger, oauth)

	// admin handlers
	adminProblemHandler := adminHandler.NewAdminProblemHandler(adminProblemStore, adminLogger, adminOauth)
	adminListHandler := adminHandler.NewAdminListHandler(adminListStore, adminLogger, adminOauth)
	adminTopicHandler := adminHandler.NewAdminTopicHandler(adminTopicStore, adminLogger, adminOauth)

	// analytics handlers
	userAnalyticsHandler := handlers.NewAnalyticsHandler(logger, oauth, analyticsStore)

	app := &Application{
		Logger:      logger,
		redisClient: redisClient,

		Oauth:      oauth,
		AdminOauth: adminOauth,

		MiddlewareHandler:     middlewareHandler,
		UserProblemHandler:    userProblemHandler,
		UserListHandler:       userListHandler,
		UserTestcaseHandler:   userTestcaseHandler,
		UserSubmissionHandler: userSubmissionHandler,
		UserTopicHandler:      userTopicHandler,

		AdminProblemHandler: adminProblemHandler,
		AdminListHandler:    adminListHandler,
		AdminTopicHandler:   adminTopicHandler,

		UserAnalyticsHandler: userAnalyticsHandler,
	}

	return app, nil
}
