package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bnquan27/Project/models"
	"github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)
	return DB.Ping()
}

// ============================================================
// Users
// ============================================================

func GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}
	err := DB.QueryRow(
		`SELECT id, email, password_hash, full_name, phone, address, role, created_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Address, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func GetUserByID(id int) (*models.User, error) {
	u := &models.User{}
	err := DB.QueryRow(
		`SELECT id, email, password_hash, full_name, phone, address, role, created_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Address, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func CreateUser(req models.RegisterRequest, passwordHash string) (*models.User, error) {
	u := &models.User{}
	err := DB.QueryRow(
		`INSERT INTO users (email, password_hash, full_name, phone, address)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, email, password_hash, full_name, phone, address, role, created_at`,
		req.Email, passwordHash, req.FullName, req.Phone, req.Address,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Address, &u.Role, &u.CreatedAt)
	return u, err
}

// ============================================================
// Sessions
// ============================================================

func CreateSession(userID int, token string, expiresAt time.Time) (*models.Session, error) {
	s := &models.Session{}
	err := DB.QueryRow(
		`INSERT INTO sessions (user_id, token, expires_at)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, token, expires_at, created_at`,
		userID, token, expiresAt,
	).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	return s, err
}

func GetSessionByToken(token string) (*models.Session, error) {
	s := &models.Session{}
	err := DB.QueryRow(
		`SELECT id, user_id, token, expires_at, created_at
		 FROM sessions WHERE token = $1 AND expires_at > NOW()`, token,
	).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func DeleteSession(token string) error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func DeleteUserSessions(userID int) error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

func UpdateUser(userID int, req models.UpdateProfileRequest) (*models.User, error) {
	u := &models.User{}
	err := DB.QueryRow(
		`UPDATE users SET full_name = $1, phone = $2, address = $3
		 WHERE id = $4
		 RETURNING id, email, password_hash, full_name, phone, address, role, created_at`,
		req.FullName, req.Phone, req.Address, userID,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Address, &u.Role, &u.CreatedAt)
	return u, err
}

func ChangePassword(userID int, newPasswordHash string) error {
	_, err := DB.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, newPasswordHash, userID)
	return err
}

// ============================================================
// Categories
// ============================================================

func GetCategories() ([]models.Category, error) {
	rows, err := DB.Query(`SELECT id, name, slug, icon, created_at FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Icon, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

// ============================================================
// Products
// ============================================================

type ProductFilter struct {
	CategoryID *int
	Search     string
	Sort       string // "price_asc", "price_desc", "newest", "name"
	Page       int
	Limit      int
	Featured   bool
}

type ProductListResult struct {
	Products   []models.Product `json:"products"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	TotalPages int              `json:"total_pages"`
}

func scanProduct(scanner interface {
	Scan(dest ...interface{}) error
}, p *models.Product) error {
	var specsJSON []byte
	var oldPrice sql.NullInt64
	var categoryID sql.NullInt64

	err := scanner.Scan(
		&p.ID, &categoryID, &p.Name, &p.Slug, &p.Description,
		&specsJSON, &p.Price, &oldPrice, pq.Array(&p.Images),
		&p.Stock, &p.Featured, &p.CreatedAt, &p.UpdatedAt,
		&p.CategoryName,
	)
	if err != nil {
		return err
	}

	if categoryID.Valid {
		cid := int(categoryID.Int64)
		p.CategoryID = &cid
	}
	if oldPrice.Valid {
		op := oldPrice.Int64
		p.OldPrice = &op
	}
	if len(specsJSON) > 0 {
		json.Unmarshal(specsJSON, &p.Specs)
	}
	if p.Specs == nil {
		p.Specs = map[string]string{}
	}
	if p.Images == nil {
		p.Images = []string{}
	}
	return nil
}

func GetProducts(filter ProductFilter) (*ProductListResult, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if filter.Featured {
		where = append(where, "featured = true")
	}
	if filter.CategoryID != nil {
		where = append(where, fmt.Sprintf("category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}
	if filter.Search != "" {
		where = append(where, fmt.Sprintf("to_tsvector('simple', name) @@ plainto_tsquery('simple', $%d)", argIdx))
		args = append(args, filter.Search)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
	if err := DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	// Sort
	orderBy := "p.created_at DESC"
	switch filter.Sort {
	case "price_asc":
		orderBy = "p.price ASC"
	case "price_desc":
		orderBy = "p.price DESC"
	case "newest":
		orderBy = "p.created_at DESC"
	case "name":
		orderBy = "p.name ASC"
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	offset := (filter.Page - 1) * filter.Limit
	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	query := fmt.Sprintf(
		`SELECT p.id, p.category_id, p.name, p.slug, p.description,
		        p.specs, p.price, p.old_price, p.images,
		        p.stock, p.featured, p.created_at, p.updated_at,
		        COALESCE(c.name, '') AS category_name
		 FROM products p
		 LEFT JOIN categories c ON c.id = p.category_id
		 WHERE %s
		 ORDER BY %s
		 LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIdx, argIdx+1,
	)
	args = append(args, filter.Limit, offset)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := scanProduct(rows, &p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if products == nil {
		products = []models.Product{}
	}

	return &ProductListResult{
		Products:   products,
		Total:      total,
		Page:       filter.Page,
		TotalPages: totalPages,
	}, rows.Err()
}

func GetProductByID(id int) (*models.Product, error) {
	row := DB.QueryRow(
		`SELECT p.id, p.category_id, p.name, p.slug, p.description,
		        p.specs, p.price, p.old_price, p.images,
		        p.stock, p.featured, p.created_at, p.updated_at,
		        COALESCE(c.name, '') AS category_name
		 FROM products p
		 LEFT JOIN categories c ON c.id = p.category_id
		 WHERE p.id = $1`, id,
	)
	var p models.Product
	if err := scanProduct(row, &p); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func GetProductBySlug(slug string) (*models.Product, error) {
	row := DB.QueryRow(
		`SELECT p.id, p.category_id, p.name, p.slug, p.description,
		        p.specs, p.price, p.old_price, p.images,
		        p.stock, p.featured, p.created_at, p.updated_at,
		        COALESCE(c.name, '') AS category_name
		 FROM products p
		 LEFT JOIN categories c ON c.id = p.category_id
		 WHERE p.slug = $1`, slug,
	)
	var p models.Product
	if err := scanProduct(row, &p); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func GetFeaturedProducts() (*ProductListResult, error) {
	f := ProductFilter{Featured: true, Limit: 20}
	return GetProducts(f)
}

func CreateProduct(req models.ProductCreateRequest) (*models.Product, error) {
	if req.Specs == nil {
		req.Specs = map[string]string{}
	}
	if req.Images == nil {
		req.Images = []string{}
	}
	specsJSON, _ := json.Marshal(req.Specs)

	row := DB.QueryRow(
		`INSERT INTO products (category_id, name, slug, description, specs, price, old_price, images, stock, featured)
		 VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7, $8, $9, $10)
		 RETURNING id, category_id, name, slug, description,
		           specs, price, old_price, images,
		           stock, featured, created_at, updated_at,
		           '' AS category_name`,
		req.CategoryID, req.Name, generateSlug(req.Name), req.Description,
		string(specsJSON), req.Price, req.OldPrice, pq.Array(req.Images), req.Stock, req.Featured,
	)
	var p models.Product
	if err := scanProduct(row, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdateProduct(id int, req models.ProductCreateRequest) (*models.Product, error) {
	if req.Specs == nil {
		req.Specs = map[string]string{}
	}
	if req.Images == nil {
		req.Images = []string{}
	}
	specsJSON, _ := json.Marshal(req.Specs)

	row := DB.QueryRow(
		`UPDATE products SET
			category_id = $1, name = $2, slug = $3, description = $4,
			specs = $5::jsonb, price = $6, old_price = $7,
			images = $8, stock = $9, featured = $10,
			updated_at = NOW()
		 WHERE id = $11
		 RETURNING id, category_id, name, slug, description,
		           specs, price, old_price, images,
		           stock, featured, created_at, updated_at,
		           '' AS category_name`,
		req.CategoryID, req.Name, generateSlug(req.Name), req.Description,
		string(specsJSON), req.Price, req.OldPrice, pq.Array(req.Images), req.Stock, req.Featured, id,
	)
	var p models.Product
	if err := scanProduct(row, &p); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func DeleteProduct(id int) error {
	res, err := DB.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func generateSlug(name string) string {
	// Simple slug: lowercase, replace spaces with hyphens, remove non-alphanumeric
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "đ", "d")
	slug = strings.ReplaceAll(slug, "/", "-")
	var result []rune
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result = append(result, r)
		}
	}
	slug = string(result)
	// Collapse multiple hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	// Append random suffix to prevent duplicates
	slug = fmt.Sprintf("%s-%d", slug, time.Now().UnixMilli()%100000)
	return slug
}

// ============================================================
// Cart
// ============================================================

func GetCart(userID int) ([]models.CartWithProduct, error) {
	rows, err := DB.Query(
		`SELECT c.id, c.user_id, c.product_id, c.quantity, c.created_at,
		        p.name, COALESCE(p.images[1], ''), p.price, p.slug, p.stock
		 FROM carts c
		 JOIN products p ON p.id = c.product_id
		 WHERE c.user_id = $1
		 ORDER BY c.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartWithProduct
	for rows.Next() {
		var ci models.CartWithProduct
		if err := rows.Scan(
			&ci.ID, &ci.UserID, &ci.ProductID, &ci.Quantity, &ci.CreatedAt,
			&ci.ProductName, &ci.ProductImage, &ci.ProductPrice, &ci.ProductSlug, &ci.Stock,
		); err != nil {
			return nil, err
		}
		items = append(items, ci)
	}
	if items == nil {
		items = []models.CartWithProduct{}
	}
	return items, rows.Err()
}

func AddToCart(userID, productID, quantity int) error {
	_, err := DB.Exec(
		`INSERT INTO carts (user_id, product_id, quantity)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, product_id) DO UPDATE
		 SET quantity = carts.quantity + $3`,
		userID, productID, quantity,
	)
	return err
}

func UpdateCartQuantity(id, userID, quantity int) error {
	res, err := DB.Exec(
		`UPDATE carts SET quantity = $1 WHERE id = $2 AND user_id = $3`,
		quantity, id, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func RemoveFromCart(id, userID int) error {
	res, err := DB.Exec(
		`DELETE FROM carts WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

// ============================================================
// Orders
// ============================================================

func CreateOrder(userID int, req models.OrderRequest, cartItems []models.CartWithProduct) (*models.Order, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check stock for all items before proceeding
	var outOfStock []string
	for _, ci := range cartItems {
		if ci.Stock < ci.Quantity {
			outOfStock = append(outOfStock, ci.ProductName)
		}
	}
	if len(outOfStock) > 0 {
		return nil, fmt.Errorf("sản phẩm đã hết hàng: %s", strings.Join(outOfStock, ", "))
	}

	// Calculate total
	var total int64
	for _, item := range cartItems {
		total += item.ProductPrice * int64(item.Quantity)
	}

	// Create order
	order := &models.Order{}
	err = tx.QueryRow(
		`INSERT INTO orders (user_id, total_amount, shipping_address, phone, note)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, status, total_amount, shipping_address, phone, note, created_at`,
		userID, total, req.ShippingAddress, req.Phone, req.Note,
	).Scan(&order.ID, &order.UserID, &order.Status, &order.TotalAmount,
		&order.ShippingAddress, &order.Phone, &order.Note, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Insert order items & decrement stock
	for _, ci := range cartItems {
		_, err = tx.Exec(
			`INSERT INTO order_items (order_id, product_id, product_name, quantity, price)
			 VALUES ($1, $2, $3, $4, $5)`,
			order.ID, ci.ProductID, ci.ProductName, ci.Quantity, ci.ProductPrice,
		)
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec(
			`UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`,
			ci.Quantity, ci.ProductID,
		)
		if err != nil {
			return nil, err
		}
	}

	// Clear cart
	_, err = tx.Exec(`DELETE FROM carts WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	return order, tx.Commit()
}

func GetUserOrders(userID int) ([]models.Order, error) {
	rows, err := DB.Query(
		`SELECT o.id, o.user_id, u.full_name, o.status, o.total_amount, o.shipping_address, o.phone, o.note, o.created_at
		 FROM orders o
		 JOIN users u ON u.id = o.user_id
		 WHERE o.user_id = $1
		 ORDER BY o.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.FullName, &o.Status, &o.TotalAmount,
			&o.ShippingAddress, &o.Phone, &o.Note, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if orders == nil {
		orders = []models.Order{}
	}
	return orders, rows.Err()
}

func GetOrderByID(id, userID int) (*models.Order, error) {
	o := &models.Order{}
	err := DB.QueryRow(
		`SELECT o.id, o.user_id, u.full_name, o.status, o.total_amount, o.shipping_address, o.phone, o.note, o.created_at
		 FROM orders o
		 JOIN users u ON u.id = o.user_id
		 WHERE o.id = $1 AND o.user_id = $2`, id, userID,
	).Scan(&o.ID, &o.UserID, &o.FullName, &o.Status, &o.TotalAmount,
		&o.ShippingAddress, &o.Phone, &o.Note, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get items
	rows, err := DB.Query(
		`SELECT id, order_id, product_id, product_name, quantity, price, created_at
		 FROM order_items WHERE order_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID,
			&item.ProductName, &item.Quantity, &item.Price, &item.CreatedAt); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}
	if o.Items == nil {
		o.Items = []models.OrderItem{}
	}
	return o, rows.Err()
}

func GetOrderByIDAdmin(id int) (*models.Order, error) {
	o := &models.Order{}
	err := DB.QueryRow(
		`SELECT o.id, o.user_id, u.full_name, o.status, o.total_amount, o.shipping_address, o.phone, o.note, o.created_at
		 FROM orders o
		 JOIN users u ON u.id = o.user_id
		 WHERE o.id = $1`, id,
	).Scan(&o.ID, &o.UserID, &o.FullName, &o.Status, &o.TotalAmount,
		&o.ShippingAddress, &o.Phone, &o.Note, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(
		`SELECT id, order_id, product_id, product_name, quantity, price, created_at
		 FROM order_items WHERE order_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID,
			&item.ProductName, &item.Quantity, &item.Price, &item.CreatedAt); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}
	if o.Items == nil {
		o.Items = []models.OrderItem{}
	}
	return o, rows.Err()
}

func GetAllOrders() ([]models.Order, error) {
	rows, err := DB.Query(
		`SELECT o.id, o.user_id, u.full_name, o.status, o.total_amount, o.shipping_address, o.phone, o.note, o.created_at
		 FROM orders o
		 JOIN users u ON u.id = o.user_id
		 ORDER BY o.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.FullName, &o.Status, &o.TotalAmount,
			&o.ShippingAddress, &o.Phone, &o.Note, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if orders == nil {
		orders = []models.Order{}
	}
	return orders, rows.Err()
}

func UpdateOrderStatus(id int, status string) error {
	res, err := DB.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func CancelOrder(orderID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Collect items first, then close rows before doing updates
	rows, err := tx.Query(
		`SELECT product_id, quantity FROM order_items WHERE order_id = $1`, orderID,
	)
	if err != nil {
		return err
	}

	type item struct{ productID, quantity int }
	var items []item
	for rows.Next() {
		var it item
		if err := rows.Scan(&it.productID, &it.quantity); err != nil {
			rows.Close()
			return err
		}
		items = append(items, it)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return err
	}

	for _, it := range items {
		_, err = tx.Exec(
			`UPDATE products SET stock = stock + $1 WHERE id = $2`, it.quantity, it.productID,
		)
		if err != nil {
			return err
		}
	}

	// Update order status
	res, err := tx.Exec(`UPDATE orders SET status = 'cancelled' WHERE id = $1`, orderID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}

	return tx.Commit()
}

// ============================================================
// Comparisons
// ============================================================

func GetOrCreateComparison(userID int) (*models.Comparison, error) {
	c := &models.Comparison{}
	err := DB.QueryRow(
		`SELECT id, user_id, created_at FROM comparisons WHERE user_id = $1`, userID,
	).Scan(&c.ID, &c.UserID, &c.CreatedAt)
	if err == sql.ErrNoRows {
		err = DB.QueryRow(
			`INSERT INTO comparisons (user_id) VALUES ($1)
			 RETURNING id, user_id, created_at`, userID,
		).Scan(&c.ID, &c.UserID, &c.CreatedAt)
	}
	return c, err
}

func AddToComparison(comparisonID, productID int) error {
	_, err := DB.Exec(
		`INSERT INTO comparison_items (comparison_id, product_id)
		 VALUES ($1, $2)
		 ON CONFLICT (comparison_id, product_id) DO NOTHING`,
		comparisonID, productID,
	)
	return err
}

func RemoveFromComparison(comparisonID, productID int) error {
	_, err := DB.Exec(
		`DELETE FROM comparison_items WHERE comparison_id = $1 AND product_id = $2`,
		comparisonID, productID,
	)
	return err
}

func GetComparisonProducts(comparisonID int) ([]models.ComparisonWithProduct, error) {
	rows, err := DB.Query(
		`SELECT ci.id, ci.comparison_id, ci.product_id,
		        p.name, COALESCE(p.images[1], ''), p.price, p.slug, p.specs,
		        COALESCE(p.category_id, 0), COALESCE(c.name, '')
		 FROM comparison_items ci
		 JOIN products p ON p.id = ci.product_id
		 LEFT JOIN categories c ON c.id = p.category_id
		 WHERE ci.comparison_id = $1`, comparisonID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ComparisonWithProduct
	for rows.Next() {
		var item models.ComparisonWithProduct
		var specsJSON []byte
		if err := rows.Scan(
			&item.ID, &item.ComparisonID, &item.ProductID,
			&item.ProductName, &item.ProductImage, &item.ProductPrice, &item.ProductSlug, &specsJSON,
			&item.CategoryID, &item.CategoryName,
		); err != nil {
			return nil, err
		}
		if len(specsJSON) > 0 {
			json.Unmarshal(specsJSON, &item.Specs)
		}
		if item.Specs == nil {
			item.Specs = map[string]string{}
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.ComparisonWithProduct{}
	}
	return items, rows.Err()
}

func GetAllProducts() ([]models.Product, error) {
	rows, err := DB.Query(
		`SELECT p.id, p.category_id, p.name, p.slug, p.description,
		        p.specs, p.price, p.old_price, p.images,
		        p.stock, p.featured, p.created_at, p.updated_at,
		        COALESCE(c.name, '') AS category_name
		 FROM products p
		 LEFT JOIN categories c ON c.id = p.category_id
		 ORDER BY p.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := scanProduct(rows, &p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if products == nil {
		products = []models.Product{}
	}
	return products, rows.Err()
}
