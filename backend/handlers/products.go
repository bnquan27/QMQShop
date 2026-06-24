package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bnquan27/QMQShop/backend/database"
	"github.com/bnquan27/QMQShop/backend/middleware"
	"github.com/bnquan27/QMQShop/backend/models"
)

// GET /api/categories
func GetCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := database.GetCategories()
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải danh mục"})
		return
	}
	middleware.JSON(w, http.StatusOK, cats)
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func parseInt64s(vals []string) []int64 {
	result := make([]int64, 0, len(vals))
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n >= 0 {
			result = append(result, n)
		}
	}
	return result
}

// GET /api/products?search=&category=&sort=&page=&limit=&brand=&cpu=&ram=&gpu=&disk=&min_price=&max_price=
// GET /api/products/featured
func GetProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Check if requesting featured products
	if strings.TrimSuffix(r.URL.Path, "/") == "/api/products/featured" {
		result, err := database.GetFeaturedProducts()
		if err != nil {
			middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải sản phẩm"})
			return
		}
		middleware.JSON(w, http.StatusOK, result)
		return
	}

	filter := database.ProductFilter{
		Search:         q.Get("search"),
		Sort:           q.Get("sort"),
		Page:           1,
		Limit:          20,
		Brands:         parseCSV(q.Get("brand")),
		CPUs:           parseCSV(q.Get("cpu")),
		RAMs:           parseCSV(q.Get("ram")),
		GPUs:           parseCSV(q.Get("gpu")),
		Disks:          parseCSV(q.Get("disk")),
		ComponentTypes: parseCSV(q.Get("component_type")),
		MinPrices:      parseInt64s(q["min_price"]),
		MaxPrices:      parseInt64s(q["max_price"]),
	}

	if p, err := strconv.Atoi(q.Get("page")); err == nil && p > 0 {
		filter.Page = p
	}
	if l, err := strconv.Atoi(q.Get("limit")); err == nil && l > 0 && l <= 100 {
		filter.Limit = l
	}

	catStr := q.Get("category")
	if catStr != "" {
		if cid, err := strconv.Atoi(catStr); err == nil {
			filter.CategoryID = &cid
		}
	}

	result, err := database.GetProducts(filter)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải sản phẩm"})
		return
	}
	middleware.JSON(w, http.StatusOK, result)
}

// GET /api/products/{id}
func GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	product, err := database.GetProductByID(id)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải sản phẩm"})
		return
	}
	if product == nil {
		middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Không tìm thấy sản phẩm"})
		return
	}
	middleware.JSON(w, http.StatusOK, product)
}

// GET /api/products/filters?category=
func GetFilterValues(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var catID *int
	if v := q.Get("category"); v != "" {
		if cid, err := strconv.Atoi(v); err == nil {
			catID = &cid
		}
	}
	opts, err := database.GetProductFilterValues(catID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải bộ lọc"})
		return
	}
	middleware.JSON(w, http.StatusOK, opts)
}

// GET /api/pc-builder/components — all buildable component groups
func GetPCBuilderComponents(w http.ResponseWriter, r *http.Request) {
	groups, err := database.GetComponentTypes()
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải linh kiện"})
		return
	}
	middleware.JSON(w, http.StatusOK, groups)
}

// GET /api/pc-builder/components/{type} — components of a specific type
func GetComponentsByType(w http.ResponseWriter, r *http.Request) {
	ct := r.PathValue("type")
	valid := false
	for _, t := range models.ComponentTypes {
		if t == ct {
			valid = true
			break
		}
	}
	if !valid {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Loại linh kiện không hợp lệ"})
		return
	}
	products, err := database.GetComponentsByType(ct)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải linh kiện"})
		return
	}
	middleware.JSON(w, http.StatusOK, products)
}
