package main

import (
	"finance-tracker/config"
	"finance-tracker/internal/auth"
	"finance-tracker/internal/transactions"
	"finance-tracker/internal/user"
	"finance-tracker/pkg/database"
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
	transactionsService := transactions.NewTransactionsService(transactionsRepository)

	/* Controllers */
	auth.NewAuthController(router, authService)
	transactions.NewTransactionsController(router, transactionsService, userRepository)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Port),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
