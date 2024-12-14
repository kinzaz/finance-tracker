package main

import (
	"finance-tracker/config"
	"finance-tracker/internal/auth"
	"finance-tracker/internal/transactions"
	"finance-tracker/internal/user"
	"finance-tracker/pkg/database"
	"finance-tracker/pkg/middleware"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := database.NewDatabase()
	conf := config.LoadConfig()
	router := http.NewServeMux()

	/* Repositories */
	userRepository := user.NewUserRepository(db)
	transactionsRepository := transactions.NewTransactionRepository(db)

	/* Services */
	authService := auth.NewAuthService(userRepository)
	userService := user.NewUserService(userRepository)
	transactionsService := transactions.NewTransactionsService(transactionsRepository, userRepository)

	/* Controllers */
	auth.NewAuthController(router, authService)
	user.NewUserController(router, userService)
	transactions.NewTransactionsController(router, transactionsService)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Port),
		Handler: middleware.CORS(router),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
