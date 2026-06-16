package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/bnquan27/Project/database"
	"github.com/bnquan27/Project/middleware"
	"github.com/bnquan27/Project/models"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/register
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	// Validate
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Vui lòng điền đầy đủ thông tin"})
		return
	}
	if len(req.Password) < 6 {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Mật khẩu phải có ít nhất 6 ký tự"})
		return
	}

	// Check existing
	existing, _ := database.GetUserByEmail(req.Email)
	if existing != nil {
		middleware.JSON(w, http.StatusConflict, models.ErrorResponse{Error: "Email đã được đăng ký"})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi xử lý mật khẩu"})
		return
	}

	user, err := database.CreateUser(req, string(hash))
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tạo tài khoản"})
		return
	}

	// Auto-login: create session
	token, _ := generateToken()
	session, err := database.CreateSession(user.ID, token, time.Now().Add(7*24*time.Hour))
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tạo phiên đăng nhập"})
		return
	}

	middleware.JSON(w, http.StatusCreated, models.AuthResponse{
		Token: session.Token,
		User:  *user,
	})
}

// POST /api/login
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	if req.Email == "" || req.Password == "" {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Vui lòng nhập email và mật khẩu"})
		return
	}

	user, err := database.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		middleware.JSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "Email hoặc mật khẩu không đúng"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		middleware.JSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "Email hoặc mật khẩu không đúng"})
		return
	}

	// Delete old sessions for this user, create new one
	database.DeleteUserSessions(user.ID)
	token, _ := generateToken()
	session, err := database.CreateSession(user.ID, token, time.Now().Add(7*24*time.Hour))
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi tạo phiên đăng nhập"})
		return
	}

	middleware.JSON(w, http.StatusOK, models.AuthResponse{
		Token: session.Token,
		User:  *user,
	})
}

// POST /api/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromRequest(r)
	if user != nil {
		database.DeleteUserSessions(user.ID)
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Đã đăng xuất"})
}

// GET /api/me
func Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromRequest(r)
	if user == nil {
		middleware.JSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "Chưa đăng nhập"})
		return
	}
	middleware.JSON(w, http.StatusOK, user)
}

// PUT /api/user/profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromRequest(r)
	if user == nil {
		middleware.JSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "Chưa đăng nhập"})
		return
	}

	var req models.UpdateProfileRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	updatedUser, err := database.UpdateUser(user.ID, req)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi cập nhật thông tin"})
		return
	}

	middleware.JSON(w, http.StatusOK, updatedUser)
}

// PUT /api/user/password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromRequest(r)
	if user == nil {
		middleware.JSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "Chưa đăng nhập"})
		return
	}

	var req models.ChangePasswordRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	// Validate
	if req.NewPassword == "" || len(req.NewPassword) < 6 {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Mật khẩu mới phải có ít nhất 6 ký tự"})
		return
	}

	// Check current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		middleware.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "Mật khẩu hiện tại không đúng"})
		return
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi xử lý mật khẩu"})
		return
	}

	if err := database.ChangePassword(user.ID, string(hash)); err != nil {
		middleware.JSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "Lỗi đổi mật khẩu"})
		return
	}

	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Đổi mật khẩu thành công"})
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
