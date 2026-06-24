// QMQSHOP — Cart Manager
const Cart = {
	_items: [],

	// Load cart from API
	async load() {
		if (!Auth.isLoggedIn()) {
			this._items = [];
			return [];
		}
		try {
			this._items = await API.getCart();
			this.updateBadge();
			return this._items;
		} catch {
			this._items = [];
			return [];
		}
	},

	// Add item
	async add(productId, quantity = 1) {
		if (!Auth.isLoggedIn()) {
			Toast.info("Vui lòng đăng nhập để thêm vào giỏ hàng");
			window.location.href = "/auth.html";
			return;
		}
		try {
			await API.addToCart(productId, quantity);
			await this.load();
			Toast.success("Đã thêm vào giỏ hàng");
		} catch (err) {
			Toast.warning(err.message || "Không thể thêm vào giỏ hàng");
		}
	},

	// Update quantity
	async update(id, quantity) {
		await API.updateCart(id, quantity);
		await this.load();
	},

	// Remove item
	async remove(id) {
		await API.removeFromCart(id);
		await this.load();
	},

	// Get items
	getItems() {
		return this._items;
	},

	// Get total count
	getCount() {
		return this._items.reduce((sum, item) => sum + item.quantity, 0);
	},

	// Get total price
	getTotal() {
		return this._items.reduce(
			(sum, item) => sum + item.product_price * item.quantity,
			0,
		);
	},

	// Update cart badge in header
	updateBadge() {
		const badge = document.getElementById("header-cart-count");
		if (badge) {
			const count = this.getCount();
			badge.textContent = count;
			badge.style.display = count > 0 ? "flex" : "none";
		}
	},

	// Format price in VND
	formatPrice(price) {
		if (!price) return "0₫";
		return price.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".") + "₫";
	},
};

// ============================================================
// Toast Notifications
// ============================================================
const Toast = {
	_container: null,

	_ensureContainer() {
		if (!this._container) {
			this._container = document.createElement("div");
			this._container.className = "toast-container";
			document.body.appendChild(this._container);
		}
		return this._container;
	},

	_show(message, type) {
		const container = this._ensureContainer();
		const toast = document.createElement("div");
		toast.className = `toast toast-${type}`;
		const icons = {
			success: "fa-circle-check",
			error: "fa-circle-xmark",
			warning: "fa-triangle-exclamation",
			info: "fa-circle-info",
		};
		toast.innerHTML = `<i class="fa-regular ${icons[type] || icons.info}"></i> ${message}`;
		container.appendChild(toast);
		setTimeout(() => {
			toast.classList.add("removing");
			setTimeout(() => toast.remove(), 300);
		}, 3000);
	},

	success(msg) {
		this._show(msg, "success");
	},
	error(msg) {
		this._show(msg, "error");
	},
	warning(msg) {
		this._show(msg, "warning");
	},
	info(msg) {
		this._show(msg, "info");
	},
};

// ============================================================
// Price formatter (global helper)
// ============================================================
function formatPrice(price) {
	return Cart.formatPrice(price);
}

// ============================================================
// Confirm dialog (Promise-based, replaces native confirm())
// ============================================================
const Confirm = {
	show({
		title,
		message,
		confirmText = "Xác nhận",
		cancelText = "KHÔNG",
		variant = "danger",
	} = {}) {
		return new Promise((resolve) => {
			const old = document.querySelector(".confirm-overlay");
			if (old) old.remove();

			const icons = {
				danger: "fa-solid fa-ban",
				warning: "fa-regular fa-eye-slash",
			};
			const overlay = document.createElement("div");
			overlay.className = "confirm-overlay";
			overlay.innerHTML = `
        <div class="confirm-dialog confirm-${variant}">
          <div class="confirm-icon confirm-icon-${variant}">
            <i class="${icons[variant] || "fa-solid fa-triangle-exclamation"}"></i>
          </div>
          <h3 class="confirm-title">${title}</h3>
          <p class="confirm-message">${message}</p>
          <div class="confirm-actions">
            <button class="btn-outline btn-sm confirm-cancel">${cancelText}</button>
            <button class="btn-primary btn-sm confirm-confirm btn-${variant === "danger" ? "danger" : "primary"}">${confirmText}</button>
          </div>
        </div>
      `;

			document.body.appendChild(overlay);
			requestAnimationFrame(() => overlay.classList.add("show"));

			const cleanup = (result) => {
				overlay.classList.remove("show");
				overlay.addEventListener("transitionend", () => overlay.remove(), {
					once: true,
				});
				setTimeout(() => {
					if (overlay.parentNode) overlay.remove();
				}, 300);
				resolve(result);
			};

			overlay
				.querySelector(".confirm-cancel")
				.addEventListener("click", () => cleanup(false));
			overlay
				.querySelector(".confirm-confirm")
				.addEventListener("click", () => cleanup(true));
			overlay.addEventListener("click", (e) => {
				if (e.target === overlay) cleanup(false);
			});

			const keyHandler = (e) => {
				if (e.key === "Escape") {
					cleanup(false);
					document.removeEventListener("keydown", keyHandler);
				}
				if (e.key === "Enter") {
					cleanup(true);
					document.removeEventListener("keydown", keyHandler);
				}
			};
			document.addEventListener("keydown", keyHandler);
		});
	},
};
