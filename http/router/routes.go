package router

import (
	"net/http"
	"time"

	adminhandler "server/http/handlers/admin_handler"
	employeehandler "server/http/handlers/employee_handler"
	userhandler "server/http/handlers/user_handler"
	"server/http/handlers/util"
	md "server/http/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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

	// Register Rate Limitter for "/v1"
	v1Router.Use(httprate.LimitByIP(10, time.Minute))

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
		r.Route("/emp", func(r chi.Router) {
			r.Post("/new", employeehandler.CreateEmp)
			r.Post("/update", employeehandler.UpdateEmp)
			r.Get("/details", employeehandler.GetEmployee)
			r.Delete("/delete", employeehandler.DeleteEmployee)
			r.Get("/net-sal", employeehandler.NetSalary)
		})
	})

	// Admin Routes
	r.Route("/admin", func(r chi.Router) {
		// Middleware
		r.Use(md.JWTMiddleware)        // Has to be a legit User
		r.Use(md.CheckAdminMiddleware) // Has to be Admin user

		// Admin Routes
		r.Get("/sal-metrics", employeehandler.GetSalaryMetricsByCountry) // Get Salary Metrics
		r.Get("/sal-avg", employeehandler.GetAvgSalaryPerJobTitle)
	})

	// Supreme Leader Route ⚡️⚡️
	r.Route("/supreme-leader", func(r chi.Router) {
		// Middleware
		r.Use(md.JWTMiddleware)           // Has to a legit User
		r.Use(md.SupremeLeaderMiddleware) // Check for Supreme Leader

		// Supreme Leader only can make or break an Admin ⚡️⚡️
		r.Post("/make-break", adminhandler.MakeBreak)
	})
}
