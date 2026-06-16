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
      Toast.info('Vui lòng đăng nhập để thêm vào giỏ hàng');
      window.location.href = '/auth.html';
      return;
    }
    await API.addToCart(productId, quantity);
    await this.load();
    Toast.success('Đã thêm vào giỏ hàng');
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
  getItems() { return this._items; },

  // Get total count
  getCount() {
    return this._items.reduce((sum, item) => sum + item.quantity, 0);
  },

  // Get total price
  getTotal() {
    return this._items.reduce((sum, item) => sum + item.product_price * item.quantity, 0);
  },

  // Update cart badge in header
  updateBadge() {
    const badge = document.getElementById('header-cart-count');
    if (badge) {
      const count = this.getCount();
      badge.textContent = count;
      badge.style.display = count > 0 ? 'flex' : 'none';
    }
  },

  // Format price in VND
  formatPrice(price) {
    if (!price) return '0₫';
    return price.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.') + '₫';
  }
};

// ============================================================
// Toast Notifications
// ============================================================
const Toast = {
  _container: null,

  _ensureContainer() {
    if (!this._container) {
      this._container = document.createElement('div');
      this._container.className = 'toast-container';
      document.body.appendChild(this._container);
    }
    return this._container;
  },

  _show(message, type) {
    const container = this._ensureContainer();
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    const icons = { success: 'fa-circle-check', error: 'fa-circle-xmark', warning: 'fa-triangle-exclamation', info: 'fa-circle-info' };
    toast.innerHTML = `<i class="fa-regular ${icons[type] || icons.info}"></i> ${message}`;
    container.appendChild(toast);
    setTimeout(() => {
      toast.classList.add('removing');
      setTimeout(() => toast.remove(), 300);
    }, 3000);
  },

  success(msg) { this._show(msg, 'success'); },
  error(msg) { this._show(msg, 'error'); },
  warning(msg) { this._show(msg, 'warning'); },
  info(msg) { this._show(msg, 'info'); }
};

// ============================================================
// Price formatter (global helper)
// ============================================================
function formatPrice(price) {
  return Cart.formatPrice(price);
}
