package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bnquan27/Project/database"
	"github.com/bnquan27/Project/middleware"
	"github.com/bnquan27/Project/models"
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

// GET /api/products?search=&category=&sort=&page=&limit=
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
		Search: q.Get("search"),
		Sort:   q.Get("sort"),
		Page:   1,
		Limit:  20,
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
