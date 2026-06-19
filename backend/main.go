package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bnquan27/QMQShop/backend/database"
	"github.com/bnquan27/QMQShop/backend/handlers"
	"github.com/bnquan27/QMQShop/backend/middleware"
)

func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func main() {
	loadEnv(".env")

	// Database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/QMQSHOP?sslmode=disable"
	}
	if err := database.InitDB(dbURL); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Connected to database")

	mux := http.NewServeMux()

	// ============================================================
	// Auth endpoints
	// ============================================================
	mux.HandleFunc("POST /api/register", handlers.Register)
	mux.HandleFunc("POST /api/login", handlers.Login)
	mux.HandleFunc("POST /api/logout", handlers.Logout)
	mux.HandleFunc("GET /api/me", handlers.Me)

	// ============================================================
	// Product & category endpoints (public)
	// ============================================================
	mux.HandleFunc("GET /api/categories", handlers.GetCategories)
	mux.HandleFunc("GET /api/products/featured", handlers.GetProducts)
	mux.HandleFunc("GET /api/products/{id}", handlers.GetProduct)
	mux.HandleFunc("GET /api/products", handlers.GetProducts)

	// ============================================================
	// Cart endpoints (require auth)
	// ============================================================
	mux.HandleFunc("GET /api/cart", middleware.RequireAuth(handlers.GetCart))
	mux.HandleFunc("POST /api/cart", middleware.RequireAuth(handlers.AddToCart))
	mux.HandleFunc("PUT /api/cart/{id}", middleware.RequireAuth(handlers.UpdateCartItem))
	mux.HandleFunc("DELETE /api/cart/{id}", middleware.RequireAuth(handlers.RemoveFromCart))

	// ============================================================
	// Order endpoints (require auth)
	// ============================================================
	mux.HandleFunc("POST /api/orders", middleware.RequireAuth(handlers.PlaceOrder))
	mux.HandleFunc("GET /api/orders", middleware.RequireAuth(handlers.GetOrders))
	mux.HandleFunc("GET /api/orders/{id}", middleware.RequireAuth(handlers.GetOrder))
	mux.HandleFunc("PUT /api/orders/{id}/cancel", middleware.RequireAuth(handlers.CancelOrder))

	// ============================================================
	// User profile endpoints (require auth)
	// ============================================================
	mux.HandleFunc("PUT /api/user/profile", middleware.RequireAuth(handlers.UpdateProfile))
	mux.HandleFunc("PUT /api/user/password", middleware.RequireAuth(handlers.ChangePassword))

	// ============================================================
	// Compare endpoints (require auth)
	// ============================================================
	mux.HandleFunc("GET /api/compare", middleware.RequireAuth(handlers.GetComparison))
	mux.HandleFunc("POST /api/compare", middleware.RequireAuth(handlers.AddToComparison))
	mux.HandleFunc("DELETE /api/compare/{productId}", middleware.RequireAuth(handlers.RemoveFromComparison))

	// ============================================================
	// Admin endpoints (require admin)
	// ============================================================
	mux.HandleFunc("GET /api/admin/products", middleware.RequireAdmin(handlers.AdminGetProducts))
	mux.HandleFunc("POST /api/admin/products", middleware.RequireAdmin(handlers.AdminCreateProduct))
	mux.HandleFunc("PUT /api/admin/products/{id}", middleware.RequireAdmin(handlers.AdminUpdateProduct))
	mux.HandleFunc("DELETE /api/admin/products/{id}", middleware.RequireAdmin(handlers.AdminDeleteProduct))
	mux.HandleFunc("GET /api/admin/orders", middleware.RequireAdmin(handlers.AdminGetOrders))
	mux.HandleFunc("GET /api/admin/orders/{id}", middleware.RequireAdmin(handlers.AdminGetOrderDetail))
	mux.HandleFunc("PUT /api/admin/orders/{id}", middleware.RequireAdmin(handlers.AdminUpdateOrderStatus))

	// ============================================================
	// Static files — serve frontend directory
	// ============================================================
	frontendDir := "../frontend"
	fs := http.FileServer(http.Dir(frontendDir))
	mux.Handle("GET /", fs)

	// ============================================================
	// Middleware stack
	// ============================================================
	handler := middleware.Logging(middleware.CORS(mux))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	log.Printf("Frontend: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
