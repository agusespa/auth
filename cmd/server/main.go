package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/agusespa/a3n/internal/database"
	"github.com/agusespa/a3n/internal/handlers"
	"github.com/agusespa/a3n/internal/helpers"
	"github.com/agusespa/a3n/internal/logger"
	"github.com/agusespa/a3n/internal/models"
	"github.com/agusespa/a3n/internal/repository"
	"github.com/agusespa/a3n/internal/service"
)

var logg logger.Logger

func corsMiddleware(next http.Handler, domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", domain)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//go:embed config/config.json
var configFile embed.FS

func main() {
	var devFlag bool
	flag.BoolVar(&devFlag, "dev", false, "enable development mode")
	flag.Parse()
	logg = logger.NewLogger(devFlag)

	encryptionKey, emailApiKey, err := helpers.GetApiKeyVars()
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to read env variables: %s", err.Error()))
	}

	dbUser, dbAddr, dbPassword, err := helpers.GetDatabaseVars()
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to read env variables: %s", err.Error()))
	}
	databaseConfig := models.Database{User: dbUser, Address: dbAddr, Password: dbPassword}
	db, err := database.ConnectDB(databaseConfig)
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to establish database connection: %s", err.Error()))
	}

	authRepository := repository.NewMySqlRepository(db)

	realmService := service.NewDefaultRealmService(authRepository, logg)
	realmEntity, err := realmService.GetRealmById(1)
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to read realm settings: %s", err.Error()))
	}

	tokenConfig := models.Token{RefreshExp: realmEntity.RefreshExp, AccessExp: realmEntity.AccessExp}

	emailSender := models.Sender{Name: realmEntity.EmailSender, Address: realmEntity.EmailAddr}
	emailConfig := models.Email{Provider: realmEntity.EmailProvider, Sender: emailSender, HardVerify: realmEntity.EmailVerify, ApiKey: emailApiKey}

	apiConfig := models.ApiConfig{Domain: realmEntity.RealmDomain, Database: databaseConfig, Token: tokenConfig, Email: emailConfig}

	emailService := service.NewDefaultEmailService(apiConfig, logg)

	apiService := service.NewDefaultApiService(authRepository, apiConfig, emailService, encryptionKey, logg)

	apiHandler := handlers.NewDefaultApiHandler(apiService, logg)

	adminHandler := handlers.NewDefaultAdminHandler(apiService, logg)

	mux := http.NewServeMux()
	// mux.HandleFunc("/api/config", apiHandler.HandleSettings)
	mux.HandleFunc("/api/register", apiHandler.HandleUserRegister)
	mux.HandleFunc("/api/login", apiHandler.HandleLogin)
	mux.HandleFunc("/api/user/email/verify", apiHandler.HandleUserEmailVerification)
	mux.HandleFunc("/api/user/email", apiHandler.HandleUserEmailChange)
	mux.HandleFunc("/api/user/password", apiHandler.HandleUserPasswordChange)
	mux.HandleFunc("/api/user", apiHandler.HandleUserData)
	mux.HandleFunc("/api/authenticate", apiHandler.HandleUserAuthentication)
	mux.HandleFunc("/api/refresh", apiHandler.HandleRefresh)
	mux.HandleFunc("/api/logout/all", apiHandler.HandleAllUserTokensRevocation)
	mux.HandleFunc("/api/logout", apiHandler.HandleTokenRevocation)

	mux.HandleFunc("/admin/dashboard", adminHandler.HandleAdminDashboard)
	mux.HandleFunc("/admin/login", adminHandler.HandleAdminLogin)

	port := os.Getenv("A3N_PORT")
	if port == "" {
		port = "3001"
	}

	logg.LogInfo(fmt.Sprintf("Listening on port %v", port))

	err = http.ListenAndServe(":"+port, corsMiddleware(mux, realmEntity.RealmDomain))
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to start HTTP server: %s", err.Error()))
	}
}
