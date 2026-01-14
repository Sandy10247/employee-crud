package router

import (
	"net/http"

	employeehandler "server/http/handlers/employee_handler"
	userhandler "server/http/handlers/user_handler"
	"server/http/handlers/util"
	md "server/http/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func InitRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(md.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	registerUtilRoutes(v1Router)
	registerUserRoutes(v1Router)

	router.Mount("/v1", v1Router)

	return router
}

func registerUtilRoutes(r chi.Router) {
	r.Get("/health", util.HandlerReady)
	r.Get("/err", util.HandleErr)
}

func registerUserRoutes(r chi.Router) {
	r.Post("/register", userhandler.HandlerCreateUser)
	r.Post("/login", userhandler.HandlerLogin)

	// Protected Routes "/v1"
	r.Route("/", func(r chi.Router) {
		// ✚ Auth Middleware
		r.Use(md.JWTMiddleware)

		// User 😊
		r.Get("/status", userhandler.CheckStatus)
		r.Get("/logout", userhandler.LogOut)

		// Employee 🤵
		r.Post("/emp/new", employeehandler.CreateEmp)
		r.Post("/emp/update", employeehandler.UpdateEmp)
		r.Get("/emp/details", employeehandler.GetEmployee)
		r.Delete("/emp/delete", employeehandler.DeleteEmployee)
	})
}
