package main

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
	"github.com/smallbatch-apps/earnsmart-api/database"
	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/routes"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

func main() {
	database.Connect()
	tbClient, err := database.CreateTigerBeetleClient()

	if err != nil {
		panic(err)
	}

	priceService := services.NewPriceService(database.DB, tbClient)
	fundService := services.NewFundService(database.DB, tbClient)
	transactionService := services.NewTransactionService(database.DB, tbClient)
	settingService := services.NewSettingService(database.DB, tbClient)
	userService := services.NewUserService(database.DB, tbClient)
	adminService := services.NewAdminService(database.DB, tbClient)
	// accountService := services.NewAccountService(database.DB, &tbClient)

	priceController := controllers.NewPriceController(priceService)
	fundController := controllers.NewFundController(fundService)
	transactionController := controllers.NewTransactionController(transactionService)
	settingController := controllers.NewSettingController(settingService)
	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)

	authedStack := middleware.CreateStack(
		middleware.LogRequest,
		middleware.RequireAuth,
		middleware.AddHeaders,
	)

	router := http.NewServeMux()
	router.HandleFunc("POST /auth/login", userController.LogIn)
	router.HandleFunc("POST /user", userController.AddUser)

	authedRouter := http.NewServeMux()
	routes.RegisterPriceRoutes(authedRouter, priceController)
	routes.RegisterFundRoutes(authedRouter, fundController)
	routes.RegisterTransactionRoutes(authedRouter, transactionController)
	routes.RegisterSettingRoutes(authedRouter, settingController)
	routes.RegisterUserRoutes(authedRouter, userController)
	routes.RegisterAdminRoutes(authedRouter, adminController)

	router.Handle("/", authedStack(authedRouter))

	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	server.ListenAndServe()
}
