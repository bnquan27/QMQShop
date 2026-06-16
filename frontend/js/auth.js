// QMQSHOP — Auth Manager

// ============================================================
// Theme Manager (light/dark toggle)
// ============================================================
const Theme = {
  KEY: 'qmq_theme',

  init() {
    const saved = localStorage.getItem(this.KEY);
    if (saved === 'light') {
      document.documentElement.setAttribute('data-theme', 'light');
    }
  },

  get() {
    return document.documentElement.getAttribute('data-theme') || 'dark';
  },

  set(theme, animate) {
    if (animate !== false) {
      document.documentElement.classList.add('theme-transitioning');
    }
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem(this.KEY, theme);
    if (animate !== false) {
      setTimeout(() => {
        document.documentElement.classList.remove('theme-transitioning');
      }, 500);
    }
  },

  toggle() {
    const next = this.get() === 'light' ? 'dark' : 'light';
    this.set(next);
    document.querySelectorAll('.theme-btn i').forEach(el => {
      el.className = next === 'light' ? 'fa-solid fa-sun' : 'fa-solid fa-moon';
    });
    return next;
  }
};

Theme.init();

const Auth = {
  _user: null,
  _listeners: [],

  // Initialize — check existing session
  async init() {
    const token = API.getToken();
    if (!token) return;
    try {
      this._user = await API.me();
      this._notify();
    } catch {
      localStorage.removeItem('qmq_token');
    }
  },

  // Login
  async login(email, password) {
    const res = await API.login(email, password);
    localStorage.setItem('qmq_token', res.token);
    this._user = res.user;
    this._notify();
    return res;
  },

  // Register
  async register(data) {
    const res = await API.register(data);
    localStorage.setItem('qmq_token', res.token);
    this._user = res.user;
    this._notify();
    return res;
  },

  // Logout
  async logout() {
    try { await API.logout(); } catch {}
    localStorage.removeItem('qmq_token');
    this._user = null;
    this._notify();
  },

  // State
  getUser() { return this._user; },
  isLoggedIn() { return !!this._user; },
  isAdmin() { return this._user && this._user.role === 'admin'; },

  // Subscribe to auth changes
  onChange(fn) {
    this._listeners.push(fn);
    return () => {
      this._listeners = this._listeners.filter(l => l !== fn);
    };
  },

  _notify() {
    this._listeners.forEach(fn => fn(this._user));
  },

  // ============================================================
  // UI Helpers
  // ============================================================

  // Update header auth buttons based on login state
  updateHeader() {
    const container = document.getElementById('auth-header-section');
    if (!container) return;

    const user = this._user;
    if (user) {
      const isLight = Theme.get() === 'light';
      container.innerHTML = `
        <button class="theme-btn" onclick="Theme.toggle(); this.classList.remove('spin'); void this.offsetWidth; this.classList.add('spin')" title="Chế độ sáng/tối">
          <i class="fa-solid ${isLight ? 'fa-sun' : 'fa-moon'}"></i>
        </button>
        <div class="user-menu">
          <button class="user-btn" onclick="event.preventDefault(); this.parentElement.classList.toggle('show')">
            <i class="fa-regular fa-user"></i>
            ${user.full_name.split(' ').pop()}
            <i class="fa-solid fa-chevron-down" style="font-size:10px;color:var(--text-muted)"></i>
          </button>
          <div class="user-dropdown">
            ${user.role === 'admin' ? '<a href="/admin/dashboard.html"><i class="fa-solid fa-gauge-high"></i> Quản trị</a>' : ''}
            <a href="/profile.html"><i class="fa-solid fa-user-gear"></i> Tài khoản</a>
            <a href="/orders.html"><i class="fa-solid fa-box"></i> Đơn hàng</a>
            <div class="dropdown-divider"></div>
            <button onclick="Auth.logout().then(() => window.location.reload())">
              <i class="fa-solid fa-right-from-bracket"></i> Đăng xuất
            </button>
          </div>
        </div>
        <a href="/cart.html" class="cart-btn" id="header-cart-btn">
          <i class="fa-solid fa-bag-shopping"></i>
          <span class="cart-count" id="header-cart-count">0</span>
        </a>
      `;
    } else {
      container.innerHTML = `
        <button class="theme-btn" onclick="Theme.toggle(); this.classList.remove('spin'); void this.offsetWidth; this.classList.add('spin')" title="Chế độ sáng/tối">
          <i class="fa-solid ${Theme.get() === 'light' ? 'fa-sun' : 'fa-moon'}"></i>
        </button>
        <a href="/auth.html" class="auth-btn">
          <i class="fa-regular fa-user"></i> Đăng nhập
        </a>
        <a href="/cart.html" class="cart-btn" id="header-cart-btn">
          <i class="fa-solid fa-bag-shopping"></i>
          <span class="cart-count" id="header-cart-count">0</span>
        </a>
      `;
    }

    // Update cart count after rendering
    if (window.Cart) {
      Cart.updateBadge();
    }
  }
};
