package main

import (
	"finance-tracker/config"
	"finance-tracker/internal/auth"
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

	/* Services */
	authService := auth.NewAuthService(userRepository)

	/* Controllers */
	auth.NewAuthController(router, authService)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Port),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
