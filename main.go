package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/smallbatch-apps/earnsmart-api/controllers"
	"github.com/smallbatch-apps/earnsmart-api/database"
	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/routes"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

func main() {
	fmt.Println("Starting server...")
	godotenv.Load()
	database.Connect()
	tbClient, err := database.CreateTigerBeetleClient()

	if err != nil {
		panic(err)
	}

	fmt.Println("Setting up services")
	services := services.NewServices(database.DB, tbClient)

	fmt.Println("Setting up controllers")
	accountController := controllers.NewAccountController(services)
	priceController := controllers.NewPriceController(services)
	fundController := controllers.NewFundController(services)
	transactionController := controllers.NewTransactionController(services)
	settingController := controllers.NewSettingController(services)
	userController := controllers.NewUserController(services)
	adminController := controllers.NewAdminController(services)
	quoteController := controllers.NewQuoteController(services)
	activityController := controllers.NewActivityController(services)
	swapController := controllers.NewSwapController(services)

	publicStack := middleware.CreateStack(
		middleware.LogRequest,
		middleware.AddHeaders, // CORS and JSON headers for public API endpoints
	)

	authedStack := middleware.CreateStack(
		middleware.LogRequest,
		middleware.RequireAuth,
		middleware.AddHeaders,
	)

	fmt.Println("Setting up routing")
	router := http.NewServeMux()
	router.Handle("POST /auth/login", publicStack(http.HandlerFunc(userController.LogIn)))
	router.Handle("POST /user", publicStack(http.HandlerFunc(userController.CreateUser)))
	router.HandleFunc("GET /seed", adminController.SeedData)
	router.HandleFunc("POST /quote", quoteController.GetQuote)
	routes.RegisterPriceRoutes(router, priceController)

	authedRouter := http.NewServeMux()
	routes.RegisterAccountRoutes(authedRouter, accountController)
	routes.RegisterFundRoutes(authedRouter, fundController)
	routes.RegisterTransactionRoutes(authedRouter, transactionController)
	routes.RegisterSettingRoutes(authedRouter, settingController)
	routes.RegisterUserRoutes(authedRouter, userController)
	routes.RegisterActivityRoutes(authedRouter, activityController)
	routes.RegisterSwapRoutes(authedRouter, swapController)

	router.Handle("/", authedStack(authedRouter))

	server := http.Server{
		Addr:    ":8090",
		Handler: middleware.AddHeaders(router),
	}

	fmt.Println("Server started at http://localhost:8090")
	if err = server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server")
	}
}
