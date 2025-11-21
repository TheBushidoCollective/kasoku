package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/thebushidocollective/kasoku/server/internal/api"
	"github.com/thebushidocollective/kasoku/server/internal/auth"
	"github.com/thebushidocollective/kasoku/server/internal/billing"
	"github.com/thebushidocollective/kasoku/server/internal/cache"
	"github.com/thebushidocollective/kasoku/server/internal/db"
	"github.com/thebushidocollective/kasoku/server/internal/jobs"
	"github.com/thebushidocollective/kasoku/server/internal/storage"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Get configuration from environment
	port := getEnv("PORT", "8080")
	dbDriver := getEnv("DB_DRIVER", "sqlite")
	dbDSN := getEnv("DB_DSN", "kasoku.db")
	jwtSecret := getEnv("JWT_SECRET", "change-me-in-production")
	storageType := getEnv("STORAGE_TYPE", "local")
	storagePath := getEnv("STORAGE_PATH", "./storage")

	// Connect to database
	database, err := db.Connect(dbDriver, dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Connected to database (driver: %s)", dbDriver)

	// Initialize storage backend
	var storageBackend storage.Storage
	switch storageType {
	case "local":
		storageBackend, err = storage.NewLocalStorage(storagePath)
		if err != nil {
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
		log.Printf("Using local storage: %s", storagePath)

	case "s3":
		s3Bucket := getEnv("S3_BUCKET", "")
		s3Region := getEnv("S3_REGION", "us-east-1")
		if s3Bucket == "" {
			log.Fatal("S3_BUCKET is required when STORAGE_TYPE=s3")
		}
		storageBackend, err = storage.NewS3Storage(nil, s3Bucket, s3Region)
		if err != nil {
			log.Fatalf("Failed to initialize S3 storage: %v", err)
		}
		log.Printf("Using S3 storage: bucket=%s, region=%s", s3Bucket, s3Region)

	default:
		log.Fatalf("Unsupported storage type: %s", storageType)
	}

	// Initialize services
	authService := auth.NewService(database, jwtSecret)
	cacheService := cache.NewService(database, storageBackend)
	jobCoordinator := jobs.NewCoordinator()

	// Initialize billing (optional - only if Stripe keys are configured)
	var billingHandler *api.BillingHandler
	stripeKey := getEnv("STRIPE_SECRET_KEY", "")
	if stripeKey != "" {
		billingService := billing.NewService(database)
		billingHandler = api.NewBillingHandler(billingService)
		log.Printf("Billing enabled with Stripe")
	}

	// Initialize OAuth config
	baseURL := getEnv("BASE_URL", "http://localhost:8080")
	oauthConfig := auth.NewOAuthConfig(baseURL)
	// Set OAuth credentials from environment
	oauthConfig.GitHub.ClientID = getEnv("GITHUB_CLIENT_ID", "")
	oauthConfig.GitHub.ClientSecret = getEnv("GITHUB_CLIENT_SECRET", "")
	oauthConfig.GitLab.ClientID = getEnv("GITLAB_CLIENT_ID", "")
	oauthConfig.GitLab.ClientSecret = getEnv("GITLAB_CLIENT_SECRET", "")

	// Initialize handlers
	authHandler := api.NewAuthHandler(authService)
	oauthHandler := api.NewOAuthHandler(authService, oauthConfig)
	cacheHandler := api.NewCacheHandler(cacheService)
	jobsHandler := api.NewJobsHandler(jobCoordinator)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-Cache-Command", "X-Cache-Created-At"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	// OAuth routes
	r.Get("/auth/github", oauthHandler.GitHubLogin)
	r.Get("/auth/gitlab", oauthHandler.GitLabLogin)
	r.Get("/auth/callback/github", oauthHandler.GitHubCallback)
	r.Get("/auth/callback/gitlab", oauthHandler.GitLabCallback)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Billing webhook (public, no auth required)
	if billingHandler != nil {
		r.Post("/billing/webhook", billingHandler.HandleWebhook)
	}

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authService.Middleware)

		// Auth endpoints
		r.Get("/auth/me", authHandler.Me)
		r.Post("/auth/tokens", authHandler.CreateToken)
		r.Get("/auth/tokens", authHandler.ListTokens)
		r.Delete("/auth/tokens/{id}", authHandler.RevokeToken)

		// Cache endpoints
		r.Get("/cache", cacheHandler.ListCache)
		r.Put("/cache/{hash}", cacheHandler.PutCache)
		r.Get("/cache/{hash}", cacheHandler.GetCache)
		r.Delete("/cache/{hash}", cacheHandler.DeleteCache)

		// Analytics endpoints
		r.Get("/analytics", cacheHandler.GetAnalytics)

		// Job coordination endpoints
		r.Post("/jobs/register/{hash}", jobsHandler.RegisterJob)
		r.Post("/jobs/complete/{hash}", jobsHandler.CompleteJob)
		r.Get("/jobs/status/{hash}", jobsHandler.GetJobStatus)
		r.Get("/jobs/wait/{hash}", jobsHandler.WaitForJob)

		// Billing endpoints (if enabled)
		if billingHandler != nil {
			r.Post("/billing/checkout", billingHandler.CreateCheckout)
			r.Post("/billing/portal", billingHandler.CreatePortal)
		}
	})

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting kasoku-server on %s", addr)
	if billingHandler != nil {
		log.Printf("Mode: Cloud with billing enabled")
	} else {
		log.Printf("Mode: Self-hosted (billing disabled)")
	}

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
