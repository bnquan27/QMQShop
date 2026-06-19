package models

import (
	"encoding/json"
	"time"
)

// ============================================================
// User
// ============================================================
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// ============================================================
// Session (token-based auth)
// ============================================================
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ============================================================
// Category
// ============================================================
type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

// ============================================================
// Product
// ============================================================
type Product struct {
	ID           int               `json:"id"`
	CategoryID   *int              `json:"category_id"`
	CategoryName string            `json:"category_name"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	Description  string            `json:"description"`
	Specs        map[string]string `json:"specs"`
	Price        int64             `json:"price"`
	OldPrice     *int64            `json:"old_price,omitempty"`
	Images       []string          `json:"images"`
	Stock        int               `json:"stock"`
	Featured     bool              `json:"featured"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func (p *Product) UnmarshalSpecs(raw []byte) error {
	return json.Unmarshal(raw, &p.Specs)
}

// ============================================================
// Cart
// ============================================================
type Cart struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

// CartWithProduct joins cart with product data
type CartWithProduct struct {
	Cart
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`
	ProductPrice int64  `json:"product_price"`
	ProductSlug  string `json:"product_slug"`
	Stock        int    `json:"stock"`
}

// ============================================================
// Order
// ============================================================
type Order struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	FullName         string    `json:"full_name"`
	Status           string    `json:"status"`
	TotalAmount      int64     `json:"total_amount"`
	ShippingAddress  string    `json:"shipping_address"`
	Phone            string    `json:"phone"`
	Note             string    `json:"note"`
	CreatedAt        time.Time `json:"created_at"`
	Items            []OrderItem `json:"items,omitempty"`
}

// ============================================================
// OrderItem
// ============================================================
type OrderItem struct {
	ID          int       `json:"id"`
	OrderID     int       `json:"order_id"`
	ProductID   int       `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Price       int64     `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

// ============================================================
// Comparison + ComparisonItem
// ============================================================
type Comparison struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ComparisonItem struct {
	ID            int       `json:"id"`
	ComparisonID  int       `json:"comparison_id"`
	ProductID     int       `json:"product_id"`
}

// ComparisonWithProduct
type ComparisonWithProduct struct {
	ComparisonItem
	ProductName  string            `json:"product_name"`
	ProductImage string            `json:"product_image"`
	ProductPrice int64             `json:"product_price"`
	ProductSlug  string            `json:"product_slug"`
	Specs        map[string]string `json:"specs"`
	CategoryID   int               `json:"category_id"`
	CategoryName string            `json:"category_name"`
}

// ============================================================
// Request / Response types
// ============================================================

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	FullName        string `json:"full_name"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	User     User   `json:"user"`
}

type CartAddRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CartUpdateRequest struct {
	Quantity int `json:"quantity"`
}

type OrderRequest struct {
	ShippingAddress string `json:"shipping_address"`
	Phone           string `json:"phone"`
	Note            string `json:"note"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CompareAddRequest struct {
	ProductID int `json:"product_id"`
}

// Admin product create/update
type ProductCreateRequest struct {
	CategoryID  *int              `json:"category_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Specs       map[string]string `json:"specs"`
	Price       int64             `json:"price"`
	OldPrice    *int64            `json:"old_price"`
	Images      []string          `json:"images"`
	Stock       int               `json:"stock"`
	Featured    bool              `json:"featured"`
}

type OrderStatusUpdate struct {
	Status string `json:"status"`
}
