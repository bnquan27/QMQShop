package handlers

import (
	"net/http"
	"strconv"

	"github.com/bnquan27/QMQShop/backend/database"
	"github.com/bnquan27/QMQShop/backend/middleware"
	"github.com/bnquan27/QMQShop/backend/models"
)

// GET /api/cart — get current user's cart
func GetCart(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)
	items, err := database.GetCart(user.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải giỏ hàng"})
		return
	}
	middleware.JSON(w, http.StatusOK, items)
}

// POST /api/cart — add item to cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	var req models.CartAddRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	// Verify product exists & has stock
	product, _ := database.GetProductByID(req.ProductID)
	if product == nil {
		middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Sản phẩm không tồn tại"})
		return
	}
	if product.Stock < 1 {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Sản phẩm đã hết hàng"})
		return
	}

	if err := database.AddToCart(user.ID, req.ProductID, req.Quantity); err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi thêm vào giỏ hàng"})
		return
	}

	// Return updated cart
	items, _ := database.GetCart(user.ID)
	middleware.JSON(w, http.StatusOK, items)
}

// PUT /api/cart/{id} — update cart item quantity
func UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	var req models.CartUpdateRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	if err := database.UpdateCartQuantity(id, user.ID, req.Quantity); err != nil {
		middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Không tìm thấy sản phẩm trong giỏ"})
		return
	}

	items, _ := database.GetCart(user.ID)
	middleware.JSON(w, http.StatusOK, items)
}

// DELETE /api/cart/{id} — remove item from cart
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	if err := database.RemoveFromCart(id, user.ID); err != nil {
		middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Không tìm thấy sản phẩm trong giỏ"})
		return
	}

	items, _ := database.GetCart(user.ID)
	middleware.JSON(w, http.StatusOK, items)
}
