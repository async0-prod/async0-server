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

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(httprate.LimitAll(100, time.Minute))
		r.Use(app.MiddlewareHandler.Cors)

		r.Route("/problems", func(r chi.Router) {
			r.Get("/list/{listID}", app.UserProblemHandler.HandlerGetListProblems)
			r.Get("/slug/{slug}", app.UserProblemHandler.HandlerGetProblemBySlug)
		})

		r.Route("/lists", func(r chi.Router) {
			r.Get("/", app.UserListHandler.HandlerGetAllLists)
		})

		r.Route("/topics", func(r chi.Router) {
			// r.Get("/", app.UserTopicHandler.HandlerGetAllTopics)
			r.Get("/", app.UserTopicHandler.HandlerGetTopics)
		})

		r.Route("/testcases", func(r chi.Router) {
			r.Get("/problem/{id}", app.UserTestcaseHandler.HandlerGetTestcaseByProblemID)
		})

		r.Route("/submissions", func(r chi.Router) {
			r.Post("/run/{id}", app.UserSubmissionHandler.HandlerRunSubmission)
			r.Post("/submit/{id}", app.UserSubmissionHandler.HandlerSubmitSubmission)
		})

	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(httprate.LimitAll(100, time.Minute))
		r.Use(app.MiddlewareHandler.Cors)

		r.Route("/problems", func(r chi.Router) {
			r.Get("/", app.AdminProblemHandler.HandlerGetAllProblems)
			r.Post("/", app.AdminProblemHandler.HandlerCreateProblem)

		})

		r.Route("/lists", func(r chi.Router) {
			r.Get("/", app.AdminListHandler.HandlerGetAllLists)
		})

		r.Route("/topics", func(r chi.Router) {
			r.Get("/", app.AdminTopicHandler.HandlerGetAllTopics)
		})
	})

	return r
}
