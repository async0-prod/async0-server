package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/grvbrk/async0_server/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(httprate.LimitAll(200, time.Minute))
	r.Use(app.MiddlewareHandler.RequestLogger)
	r.Use(app.MiddlewareHandler.Security)

	r.Route("/auth", func(r chi.Router) {

		r.Use(httprate.LimitAll(100, time.Minute))

		// Auth routes without CORS
		r.Get("/google/login", app.Oauth.Login)
		r.Get("/google/logout", app.Oauth.Logout)
		r.Get("/google/callback", app.Oauth.Callback)

		r.Get("/admin/google/login", app.AdminOauth.Login)
		r.Get("/admin/google/logout", app.AdminOauth.Logout)
		r.Get("/admin/google/callback", app.AdminOauth.Callback)

		// Auth routes with CORS
		r.Group(func(r chi.Router) {
			r.Use(app.MiddlewareHandler.Cors)
			r.Get("/user", app.Oauth.AuthUser)
			r.Get("/admin", app.AdminOauth.AuthAdmin)

		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(httprate.LimitAll(100, time.Minute))
		r.Use(app.MiddlewareHandler.Cors)

		r.Route("/problems", func(r chi.Router) {
			r.With(app.MiddlewareHandler.SoftAuthenticate).
				Get("/table/{listID}", app.UserProblemHandler.HandlerGetTanstackTableProblems)
			r.Get("/{slug}", app.UserProblemHandler.HandlerGetProblemBySlug)
		})

		r.Route("/solutions", func(r chi.Router) {
			r.Get("/problem/{id}", app.UserSolutionHandler.HandlerGetSolutionsByProblemID)
		})

		r.Route("/lists", func(r chi.Router) {
			r.Get("/", app.UserListHandler.HandlerGetAllLists)
		})

		r.Route("/topics", func(r chi.Router) {
			r.Get("/", app.UserTopicHandler.HandlerGetAllTopics)
			r.Get("/list/{listID}", app.UserTopicHandler.HandlerGetTopicsByListID)
		})

		r.Route("/testcases", func(r chi.Router) {
			r.Get("/problem/{id}", app.UserTestcaseHandler.HandlerGetTestcasesByProblemID)
		})

		r.Route("/submissions", func(r chi.Router) {
			r.Use(app.MiddlewareHandler.Authenticate)

			r.Get("/problem/{problemID}", app.UserSubmissionHandler.HandlerGetSubmissionsByProblemID)

			r.Post("/run", app.UserSubmissionHandler.HandlerRunSubmission)
			r.Post("/submit/{id}", app.UserSubmissionHandler.HandlerSubmitSubmission)
		})

		r.Route("/analytics", func(r chi.Router) {
			r.Use(app.MiddlewareHandler.SoftAuthenticate)
			r.Get("/list/{listID}", app.UserAnalyticsHandler.HandlerGetCardAnalyticsByListID)
		})

	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(httprate.LimitAll(100, time.Minute))
		r.Use(app.MiddlewareHandler.Cors)
		// r.Use(app.MiddlewareHandler.AuthenticateAdmin)

		r.Route("/problems", func(r chi.Router) {
			r.Get("/", app.AdminProblemHandler.HandlerGetAllProblems)
			r.Get("/{id}", app.AdminProblemHandler.HandlerGetProblemByID)
			r.Post("/", app.AdminProblemHandler.HandlerCreateProblem)
			r.Put("/{id}", app.AdminProblemHandler.HandlerUpdateProblem)
		})

		r.Route("/lists", func(r chi.Router) {
			r.Get("/", app.AdminListHandler.HandlerGetAllLists)
			r.Get("/problem/{id}", app.AdminListHandler.HandlerGetListsByProblemID)
		})

		r.Route("/topics", func(r chi.Router) {
			r.Get("/", app.AdminTopicHandler.HandlerGetAllTopics)
			r.Get("/problem/{id}", app.AdminTopicHandler.HandlerGetTopicsByProblemID)
		})

		r.Route("/testcases", func(r chi.Router) {
			r.Get("/problem/{id}", app.AdminTestcaseHandler.HandlerGetTestcasesByProblemID)
		})

		r.Route("/solutions", func(r chi.Router) {
			r.Get("/problem/{id}", app.AdminSolutionHandler.HandlerGetSolutionsByProblemID)
		})
	})

	return r
}
