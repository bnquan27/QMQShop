package handlers

import (
	"net/http"
	"strconv"

	"github.com/bnquan27/Project/database"
	"github.com/bnquan27/Project/middleware"
	"github.com/bnquan27/Project/models"
)

// POST /api/orders — place order from cart contents
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	var req models.OrderRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	// Validate
	if req.ShippingAddress == "" {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Vui lòng nhập địa chỉ giao hàng"})
		return
	}
	if req.Phone == "" {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Vui lòng nhập số điện thoại"})
		return
	}

	// Get cart
	cartItems, err := database.GetCart(user.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải giỏ hàng"})
		return
	}
	if len(cartItems) == 0 {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Giỏ hàng trống"})
		return
	}

	order, err := database.CreateOrder(user.ID, req, cartItems)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi đặt hàng"})
		return
	}

	middleware.JSON(w, http.StatusCreated, order)
}

// GET /api/orders — get current user's order history
func GetOrders(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)
	orders, err := database.GetUserOrders(user.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải đơn hàng"})
		return
	}
	middleware.JSON(w, http.StatusOK, orders)
}

// GET /api/orders/{id} — get order detail with items
func GetOrder(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	order, err := database.GetOrderByID(id, user.ID)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải đơn hàng"})
		return
	}
	if order == nil {
		middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Không tìm thấy đơn hàng"})
		return
	}
	middleware.JSON(w, http.StatusOK, order)
}
