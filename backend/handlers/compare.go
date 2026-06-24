package handlers

import (
	"net/http"
	"strconv"

	"github.com/bnquan27/QMQShop/backend/database"
	"github.com/bnquan27/QMQShop/backend/middleware"
	"github.com/bnquan27/QMQShop/backend/models"
)

// GET /api/compare — get user's comparison list with products
func GetComparison(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	comp, err := database.GetOrCreateComparison(user.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải danh sách so sánh"})
		return
	}

	items, err := database.GetComparisonProducts(comp.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải sản phẩm so sánh"})
		return
	}

	middleware.JSON(w, http.StatusOK, map[string]interface{}{
		"comparison": comp,
		"products":   items,
	})
}

// POST /api/compare — add product to comparison
func AddToComparison(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	var req models.CompareAddRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	// Verify product exists
	product, _ := database.GetProductByID(req.ProductID)
	if product == nil {
		middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Sản phẩm không tồn tại"})
		return
	}

	comp, err := database.GetOrCreateComparison(user.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tạo danh sách so sánh"})
		return
	}

	current, _ := database.GetComparisonProducts(comp.ID)
	if len(current) >= 3 {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Chỉ được so sánh tối đa 3 sản phẩm"})
		return
	}

	if len(current) > 0 {
		existingCatID := current[0].CategoryID
		if product.CategoryID != nil && *product.CategoryID != existingCatID {
			middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{
				Error: "Chỉ có thể so sánh sản phẩm trong cùng danh mục",
			})
			return
		}
		// For component-type products, enforce same component type
		existingCT := current[0].ComponentType
		if existingCT != "" && product.ComponentType != "" && product.ComponentType != existingCT {
			middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{
				Error: "Chỉ có thể so sánh linh kiện cùng loại",
			})
			return
		}
	}

	if err := database.AddToComparison(comp.ID, req.ProductID); err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi thêm vào so sánh"})
		return
	}

	items, _ := database.GetComparisonProducts(comp.ID)
	middleware.JSON(w, http.StatusOK, items)
}

// DELETE /api/compare — clear all products from comparison
func ClearComparison(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	comp, err := database.GetOrCreateComparison(user.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải danh sách so sánh"})
		return
	}

	if err := database.ClearComparison(comp.ID); err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi xóa danh sách so sánh"})
		return
	}

	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Đã xóa danh sách so sánh"})
}

// DELETE /api/compare/{productId} — remove product from comparison
func RemoveFromComparison(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	pidStr := r.PathValue("productId")
	productID, err := strconv.Atoi(pidStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	comp, err := database.GetOrCreateComparison(user.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải danh sách so sánh"})
		return
	}

	if err := database.RemoveFromComparison(comp.ID, productID); err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi xóa khỏi so sánh"})
		return
	}

	items, _ := database.GetComparisonProducts(comp.ID)
	middleware.JSON(w, http.StatusOK, items)
}
