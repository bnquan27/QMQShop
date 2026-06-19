<div align="center">
  <h1>🛒 QMQSHOP</h1>
  <p><strong>Laptop, PC Gaming & Linh Kiện Chính Hãng</strong></p>
  <p>
    <img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go" alt="Go" />
    <img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql" alt="PostgreSQL" />
    <img src="https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5" alt="HTML5" />
    <img src="https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3" alt="CSS3" />
    <img src="https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript" alt="JavaScript" />
  </p>
  <p>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="License" />
  </p>
</div>

---

## 📋 Table of Contents

- [✨ Features](#-features)
- [🛠 Tech Stack](#-tech-stack)
- [📁 Project Structure](#-project-structure)
- [⚙️ Prerequisites](#️-prerequisites)
- [🚀 Installation](#-installation)
- [💻 Usage](#-usage)
- [🔐 Admin Access](#-admin-access)
- [📸 Screenshots](#-screenshots)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🏪 **Product Catalog** | Browse laptops, PC gaming, PC văn phòng & linh kiện |
| 🔍 **Live Search** | Real-time product search with dropdown suggestions |
| 🛒 **Shopping Cart** | Add, update, remove items with persistent cart |
| ⚖️ **Compare Products** | Side-by-side comparison (max 3, same category) |
| 👤 **User Accounts** | Register, login, profile management |
| 🔐 **Admin Panel** | Manage products, orders, status updates |
| 📦 **Order Management** | Place orders, track status (pending → delivered) |
| 🌙 **Dark / Light Mode** | Theme toggle with smooth transitions |
| 📱 **Responsive** | Optimized for desktop, tablet & mobile |

---

## 🛠 Tech Stack

```
Frontend         │  Backend           │  Database
─────────────────┼────────────────────┼─────────────────
HTML5 + CSS3     │  Go 1.26          │  PostgreSQL 16
Vanilla JS       │  net/http (mux)   │  pgAdmin
Font Awesome 6   │  bcrypt + JWT     │
Google Fonts     │  CORS middleware  │
(Anton + DM Sans)│                    │
```

---

## 📁 Project Structure

```
Project/
├── backend/
│   ├── main.go                 # Entry point, routes
│   ├── go.mod / go.sum
│   ├── database/
│   │   └── db.go               # SQL queries & DB layer
│   ├── handlers/
│   │   ├── auth.go             # Login, register, profile
│   │   ├── products.go         # Product CRUD
│   │   ├── cart.go             # Cart operations
│   │   ├── compare.go          # Product comparison
│   │   ├── orders.go           # Order placement
│   │   └── admin.go            # Admin endpoints
│   ├── middleware/
│   │   └── middleware.go       # Auth, CORS, logging
│   └── models/
│       └── models.go           # Data types
│
├── frontend/
│   ├── index.html              # Homepage
│   ├── products.html           # Product listing
│   ├── product.html            # Product detail
│   ├── cart.html               # Shopping cart
│   ├── orders.html             # User orders
│   ├── auth.html               # Login / Register
│   ├── profile.html            # Account settings
│   ├── compare.html            # Compare products
│   ├── js/                     # JavaScript modules
│   ├── css/                    # Stylesheets
│   └── admin/                  # Admin dashboard
│       ├── dashboard.html
│       ├── products.html
│       └── orders.html
│
└── README.md
```

---

## ⚙️ Prerequisites

- [Go](https://go.dev/doc/install) 1.26+
- PostgreSQL 16+
- A running PostgreSQL database named `QMQSHOP`

---

## 🚀 Installation

### 1️⃣ Database Setup

Create the database and run the schema:

```sql
CREATE DATABASE "QMQSHOP";
```

Make sure your PostgreSQL connection string matches the one expected by the app:

```
postgres://postgres:postgres@localhost:5432/QMQSHOP?sslmode=disable
```

### 3️⃣ Run the Server

```bash
cd backend
go run .
```

The server starts on **`http://localhost:8080`**.

---

## 💻 Usage

| Page | URL | Description |
|------|-----|-------------|
| 🏠 Home | `/` | Featured products & categories |
| 🛍 Products | `/products.html` | Browse & search all products |
| 🔍 Product | `/product.html?id=N` | Product detail |
| 🛒 Cart | `/cart.html` | Shopping cart |
| ⚖️ Compare | `/compare.html` | Compare products |
| 👤 Profile | `/profile.html` | Edit account & password |
| 📦 Orders | `/orders.html` | Order history |
| 🔑 Auth | `/auth.html` | Login / Register |
| ⚙️ Admin | `/admin/dashboard.html` | Manage products & orders |

> ℹ️ Click **"Tất cả"** in the nav or browse categories to explore products.

---

## 🔐 Admin Access

| Credential | Value |
|------------|-------|
| **Email** | `admin@qmqshop.com` |
| **Password** | `123456` |

---

## 📸 Screenshots

<div align="center">

### 🏠 Trang Chủ
<img src="screenshots/Main_page.png" alt="Main Page" width="820" style="border-radius: 8px; border: 1px solid #333;" />
<p><em>Hero slider, danh mục sản phẩm và sản phẩm nổi bật</em></p>

<br />

### 🛍️ Danh Sách Sản Phẩm
<img src="screenshots/Products_grid.png" alt="Products Grid" width="820" style="border-radius: 8px; border: 1px solid #333;" />
<p><em>Lưới sản phẩm với tìm kiếm và phân loại theo danh mục</em></p>

<br />

| 📦 Quản Lý Đơn Hàng | ⚙️ Admin Dashboard |
|:---:|:---:|
| <img src="screenshots/Orders.png" alt="Orders" width="390" style="border-radius: 8px; border: 1px solid #333;" /> | <img src="screenshots/Admin_dashboard.png" alt="Admin Dashboard" width="390" style="border-radius: 8px; border: 1px solid #333;" /> |
| <em>Theo dõi trạng thái đơn hàng</em> | <em>Quản lý sản phẩm & đơn hàng</em> |

</div>

---

<div align="center">
  <sub>Built with ❤️ using Go & vanilla JS — QMQSHOP SINCE 2026</sub>
</div>
