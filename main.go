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
	godotenv.Load()
	database.Connect()
	tbClient, err := database.CreateTigerBeetleClient()

	if err != nil {
		panic(err)
	}

	fmt.Println("Setting up services")

	accountService := services.NewAccountService(database.DB, tbClient)
	priceService := services.NewPriceService(database.DB, tbClient)
	fundService := services.NewFundService(database.DB, tbClient)
	transactionService := services.NewTransactionService(database.DB, tbClient)
	settingService := services.NewSettingService(database.DB, tbClient)
	userService := services.NewUserService(database.DB, tbClient)
	adminService := services.NewAdminService(database.DB, tbClient)
	activityService := services.NewActivityService(database.DB, tbClient)

	fmt.Println("Setting up controllers")
	accountController := controllers.NewAccountController(accountService)
	priceController := controllers.NewPriceController(priceService)
	fundController := controllers.NewFundController(fundService)
	transactionController := controllers.NewTransactionController(transactionService)
	settingController := controllers.NewSettingController(settingService)
	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)
	quoteController := controllers.NewQuoteController(priceService)
	activityController := controllers.NewActivityController(activityService)

	authedStack := middleware.CreateStack(
		middleware.LogRequest,
		middleware.RequireAuth,
		middleware.AddHeaders,
	)

	fmt.Println("Setting up routing")
	router := http.NewServeMux()
	router.HandleFunc("POST /auth/login", userController.LogIn)
	router.HandleFunc("POST /user", userController.AddUser)
	router.HandleFunc("GET /seed", adminController.SeedData)
	router.HandleFunc("POST /quote", quoteController.GetQuote)

	authedRouter := http.NewServeMux()
	routes.RegisterAccountRoutes(authedRouter, accountController)
	routes.RegisterPriceRoutes(authedRouter, priceController)
	routes.RegisterFundRoutes(authedRouter, fundController)
	routes.RegisterTransactionRoutes(authedRouter, transactionController)
	routes.RegisterSettingRoutes(authedRouter, settingController)
	routes.RegisterUserRoutes(authedRouter, userController)

	router.Handle("/", authedStack(authedRouter))

	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	fmt.Println("Server started at http://localhost:8090")
	if err = server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server")
	}
}
