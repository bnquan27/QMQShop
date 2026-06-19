package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bnquan27/QMQShop/backend/database"
	"github.com/bnquan27/QMQShop/backend/middleware"
	"github.com/bnquan27/QMQShop/backend/models"
)

// GET /api/admin/products — all products for admin management
func AdminGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := database.GetAllProducts()
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải sản phẩm"})
		return
	}
	middleware.JSON(w, http.StatusOK, products)
}

// POST /api/admin/products — create new product
func AdminCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.ProductCreateRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	if req.Name == "" {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Tên sản phẩm không được để trống"})
		return
	}
	if req.Price <= 0 {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Giá sản phẩm không hợp lệ"})
		return
	}

	product, err := database.CreateProduct(req)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tạo sản phẩm"})
		return
	}
	middleware.JSON(w, http.StatusCreated, product)
}

// PUT /api/admin/products/{id} — update product
func AdminUpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	var req models.ProductCreateRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	product, err := database.UpdateProduct(id, req)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi cập nhật sản phẩm"})
		return
	}
	if product == nil {
		middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Không tìm thấy sản phẩm"})
		return
	}
	middleware.JSON(w, http.StatusOK, product)
}

// DELETE /api/admin/products/{id} — delete product
func AdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	if err := database.DeleteProduct(id); err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			middleware.JSON(w, http.StatusConflict, models.ErrorResponse{Error: "Không thể xóa sản phẩm đã có trong đơn hàng"})
		} else if err.Error() == "not found" {
			middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Không tìm thấy sản phẩm"})
		} else {
			middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi xóa sản phẩm"})
		}
		return
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Đã xóa sản phẩm"})
}

// GET /api/admin/orders — all orders for admin management
func AdminGetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := database.GetAllOrders()
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tải đơn hàng"})
		return
	}
	middleware.JSON(w, http.StatusOK, orders)
}

func AdminGetOrderDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	order, err := database.GetOrderByIDAdmin(id)
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

// PUT /api/admin/orders/{id} — update order status
func AdminUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "ID không hợp lệ"})
		return
	}

	var req models.OrderStatusUpdate
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	validStatuses := map[string]bool{
		"pending": true, "confirmed": true,
		"shipping": true, "delivered": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Trạng thái không hợp lệ"})
		return
	}

	if err := database.UpdateOrderStatus(id, req.Status); err != nil {
		middleware.JSON(w, http.StatusNotFound, models.ErrorResponse{Error: "Không tìm thấy đơn hàng"})
		return
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Đã cập nhật trạng thái"})
}
