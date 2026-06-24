// QMQSHOP — API Client
const API = {
  baseURL: '',

  // Get auth token
  getToken() {
    return localStorage.getItem('qmq_token');
  },

  // Generic request
  async request(method, path, body = null) {
    const headers = { 'Content-Type': 'application/json' };
    const token = this.getToken();
    if (token) {
      headers['Authorization'] = 'Bearer ' + token;
    }

    const opts = { method, headers };
    if (body !== null) {
      opts.body = JSON.stringify(body);
    }

    const res = await fetch(this.baseURL + path, opts);
    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || 'Yêu cầu thất bại');
    }
    return data;
  },

  // Auth
  login(email, password) {
    return this.request('POST', '/api/login', { email, password });
  },
  register(data) {
    return this.request('POST', '/api/register', data);
  },
  logout() {
    return this.request('POST', '/api/logout');
  },
  me() {
    return this.request('GET', '/api/me');
  },

  // Products
  getProducts(params = {}) {
    const q = new URLSearchParams();
    if (params.search) q.set('search', params.search);
    if (params.category) q.set('category', params.category);
    if (params.sort) q.set('sort', params.sort);
    if (params.page) q.set('page', params.page);
    if (params.limit) q.set('limit', params.limit);
    if (params.brand) q.set('brand', params.brand);
    if (params.component_type) q.set('component_type', params.component_type);
    if (params.cpu) q.set('cpu', params.cpu);
    if (params.ram) q.set('ram', params.ram);
    if (params.gpu) q.set('gpu', params.gpu);
    if (params.disk) q.set('disk', params.disk);
    // Support both single value and array for price ranges
    // Use explicit filter to preserve 0 values (e.g. min_price=0 for "Dưới 5 triệu")
    [params.min_price].flat().filter(v => v !== null && v !== undefined && v !== '').forEach(v => q.append('min_price', v));
    [params.max_price].flat().filter(v => v !== null && v !== undefined && v !== '').forEach(v => q.append('max_price', v));
    const qs = q.toString();
    return this.request('GET', '/api/products' + (qs ? '?' + qs : ''));
  },
  getFeatured() {
    return this.request('GET', '/api/products/featured');
  },
  getProduct(id) {
    return this.request('GET', '/api/products/' + id);
  },
  getCategories() {
    return this.request('GET', '/api/categories');
  },
  getFilterOptions(params = {}) {
    const q = new URLSearchParams();
    if (params.category) q.set('category', params.category);
    const qs = q.toString();
    return this.request('GET', '/api/products/filters' + (qs ? '?' + qs : ''));
  },

  // Cart
  getCart() {
    return this.request('GET', '/api/cart');
  },
  addToCart(productId, quantity = 1) {
    return this.request('POST', '/api/cart', { product_id: productId, quantity });
  },
  updateCart(id, quantity) {
    return this.request('PUT', '/api/cart/' + id, { quantity });
  },
  removeFromCart(id) {
    return this.request('DELETE', '/api/cart/' + id);
  },

  // Orders
  placeOrder(data) {
    return this.request('POST', '/api/orders', data);
  },
  getOrders() {
    return this.request('GET', '/api/orders');
  },
  getOrder(id) {
    return this.request('GET', '/api/orders/' + id);
  },
  cancelOrder(id) {
    return this.request('PUT', '/api/orders/' + id + '/cancel');
  },

  // User profile
  updateProfile(data) {
    return this.request('PUT', '/api/user/profile', data);
  },
  changePassword(data) {
    return this.request('PUT', '/api/user/password', data);
  },

  // AI Chat
  chat(message) {
    return this.request('POST', '/api/chat', { message });
  },

  // Compare
  getCompare() {
    return this.request('GET', '/api/compare');
  },
  addToCompare(productId) {
    return this.request('POST', '/api/compare', { product_id: productId });
  },
  removeFromCompare(productId) {
    return this.request('DELETE', '/api/compare/' + productId);
  },
  clearCompare() {
    return this.request('DELETE', '/api/compare');
  },

  // Admin
  adminGetProducts() {
    return this.request('GET', '/api/admin/products');
  },
  adminCreateProduct(data) {
    return this.request('POST', '/api/admin/products', data);
  },
  adminUpdateProduct(id, data) {
    return this.request('PUT', '/api/admin/products/' + id, data);
  },
  adminDeleteProduct(id) {
    return this.request('DELETE', '/api/admin/products/' + id);
  },
  adminToggleProductHidden(id) {
    return this.request('PUT', '/api/admin/products/' + id + '/toggle-hidden');
  },
  adminGetOrders() {
    return this.request('GET', '/api/admin/orders');
  },
  adminUpdateOrderStatus(id, status) {
    return this.request('PUT', '/api/admin/orders/' + id, { status });
  },
  adminGetOrderDetail(id) {
    return this.request('GET', '/api/admin/orders/' + id);
  },

  // Admin filter options
  adminGetFilterOptions() {
    return this.request('GET', '/api/admin/filter-options');
  },
  adminCreateFilterOption(data) {
    return this.request('POST', '/api/admin/filter-options', data);
  },
  adminUpdateFilterOption(id, data) {
    return this.request('PUT', '/api/admin/filter-options/' + id, data);
  },
  adminDeleteFilterOption(id) {
    return this.request('DELETE', '/api/admin/filter-options/' + id);
  }
};
