// QMQSHOP — Compare Manager
const Compare = {
  _items: [],

  // Load compare list from API
  async load() {
    if (!Auth.isLoggedIn()) {
      this._items = [];
      return [];
    }
    try {
      const data = await API.getCompare();
      this._items = data.products || [];
      return this._items;
    } catch {
      this._items = [];
      return [];
    }
  },

  MAX_COMPARE: 3,

  // Add product to compare
  async add(productId) {
    if (!Auth.isLoggedIn()) {
      Toast.info('Vui lòng đăng nhập để so sánh sản phẩm');
      return;
    }
    if (this._items.length >= this.MAX_COMPARE) {
      Toast.warning('Chỉ được so sánh tối đa ' + this.MAX_COMPARE + ' sản phẩm');
      return;
    }
    try {
      await API.addToCompare(productId);
      await this.load();
    } catch (err) {
      Toast.error(err.message);
    }
  },

  // Remove product from compare
  async remove(productId) {
    try {
      await API.removeFromCompare(productId);
      await this.load();
      Toast.success('Đã xoá khỏi so sánh');
    } catch (err) {
      Toast.error(err.message);
    }
  },

  getItems() { return this._items; },
  getCount() { return this._items.length; },

  // ============================================================
  // Floating Compare Modal
  // ============================================================

  // Open the floating compare modal (optionally add a product first)
  async open(productId) {
    if (!Auth.isLoggedIn()) {
      Toast.info('Vui lòng đăng nhập để so sánh sản phẩm');
      window.location.href = '/auth.html';
      return;
    }

    this._ensureModal();

    // Load current compare items
    await this.load();

    if (productId) {
      const isInList = this._items.some(p => p.product_id === productId);
      if (!isInList) {
        if (this._items.length >= this.MAX_COMPARE) {
          Toast.warning('Chỉ được so sánh tối đa ' + this.MAX_COMPARE + ' sản phẩm');
        } else {
          try {
            await API.addToCompare(productId);
            await this.load();
          } catch (err) {
            Toast.error(err.message);
          }
        }
      }
    }

    // Show modal
    const modal = document.getElementById('compare-modal');
    modal.style.display = 'flex';
    document.body.style.overflow = 'hidden';

    // Render both panels
    this._renderModalCompare();
    await this._renderProductSelector();
  },

  // Close the floating compare modal
  close() {
    const modal = document.getElementById('compare-modal');
    if (modal) modal.style.display = 'none';
    document.body.style.overflow = '';
  },

  // Create modal DOM if it doesn't exist
  _ensureModal() {
    if (document.getElementById('compare-modal')) return;

    const div = document.createElement('div');
    div.id = 'compare-modal';
    div.className = 'compare-modal';
    div.style.display = 'none';
    div.innerHTML = `
      <div class="compare-modal-backdrop" onclick="Compare.close()"></div>
      <div class="compare-modal-content">
        <div class="compare-modal-header">
          <h3><i class="fa-solid fa-scale-balanced" style="color:var(--accent);"></i> SO SÁNH SẢN PHẨM</h3>
          <button class="compare-modal-close" onclick="Compare.close()">&times;</button>
        </div>
        <div class="compare-modal-body">
          <!-- Left: Comparison table -->
          <div class="cms-compare-wrap" id="cms-compare-wrap">
            <div id="cms-empty" class="cms-empty" style="display:flex;">
              <i class="fa-regular fa-scale-balanced"></i>
              <span>Chọn sản phẩm từ bên phải để so sánh</span>
            </div>
            <div id="cms-table-wrap" style="display:none;overflow-x:auto;">
              <div class="cms-compare-category"><i class="fa-solid fa-layer-group"></i> <span id="cms-cat-label"></span></div>
              <div class="compare-grid" id="cms-compare-table"></div>
            </div>
          </div>
          <!-- Right: Product selector -->
          <div class="cms-selector">
            <div class="cms-selector-header">
              <i class="fa-solid fa-search" style="color:var(--text-muted);font-size:14px;"></i>
              <input type="text" id="cms-search-input" placeholder="Tìm sản phẩm..." oninput="Compare.searchProducts(this.value)" />
            </div>
            <div class="cms-category-notice" id="cms-category-notice" style="display:none;"></div>
            <div class="cms-products" id="cms-products">
              <div class="loading-center" style="padding:20px;"><div class="spinner"></div></div>
            </div>
          </div>
        </div>
      </div>
    `;
    document.body.appendChild(div);
  },

  // Render the comparison table inside the modal
  _renderModalCompare() {
    const products = this._items;

    if (products.length === 0) {
      document.getElementById('cms-empty').style.display = 'flex';
      document.getElementById('cms-table-wrap').style.display = 'none';
      return;
    }

    document.getElementById('cms-empty').style.display = 'none';
    document.getElementById('cms-table-wrap').style.display = 'block';

    // Show category label
    const catName = products[0].category_name;
    const catLabel = document.getElementById('cms-cat-label');
    if (catLabel) catLabel.textContent = catName || '';

    // Collect all spec keys across all products
    const allKeys = new Set();
    products.forEach(p => {
      if (p.specs) Object.keys(p.specs).forEach(k => allKeys.add(k));
    });

    const specKeys = ['Hình ảnh', 'Tên sản phẩm', 'Giá', ...Array.from(allKeys)];
    const table = document.getElementById('cms-compare-table');

    table.style.gridTemplateColumns = `200px repeat(${products.length}, minmax(240px, 1fr))`;
    table.innerHTML = specKeys.map(key => {
      const label = `<div class="compare-label">${key}</div>`;
      const cells = products.map(p => {
        if (key === 'Hình ảnh') {
          return `<div class="compare-cell" style="padding:0;justify-content:center;"><img src="${p.product_image || ''}" alt="${p.product_name}" /></div>`;
        }
        if (key === 'Tên sản phẩm') {
          return `<div class="compare-cell" style="flex-direction:column;align-items:flex-start;gap:4px;">
            <a href="/product.html?id=${p.product_id}" target="_blank" style="font-weight:600;color:var(--text-primary);font-size:13px;">${p.product_name}</a>
            <button class="btn-ghost" style="color:#dc2626;font-size:11px;padding:2px 0;" onclick="Compare.remove(${p.product_id}).then(() => { Compare._renderModalCompare(); Compare._renderProductSelector(); })"><i class="fa-regular fa-trash-can"></i> Xoá</button>
          </div>`;
        }
        if (key === 'Giá') {
          return `<div class="compare-cell"><span class="price-current" style="font-size:16px;">${formatPrice(p.product_price)}</span></div>`;
        }
        return `<div class="compare-cell">${p.specs && p.specs[key] || '—'}</div>`;
      });
      return label + cells.join('');
    }).join('');
  },

  // Render the product selector panel (right side)
  async _renderProductSelector(searchTerm) {
    const container = document.getElementById('cms-products');
    container.innerHTML = '<div class="loading-center" style="padding:20px;"><div class="spinner"></div></div>';

    try {
      // Determine category from current compare items
      const catId = this._items.length > 0 ? this._items[0].category_id : null;
      const catName = this._items.length > 0 ? this._items[0].category_name : null;

      // Show category notice
      const notice = document.getElementById('cms-category-notice');
      if (catId) {
        notice.style.display = 'flex';
        notice.innerHTML = '<i class="fa-solid fa-layer-group"></i> Danh mục: <strong>' + catName + '</strong>';
      } else {
        notice.style.display = 'none';
      }

      const params = { limit: 30 };
      if (catId) params.category = catId;
      if (searchTerm && searchTerm.trim()) params.search = searchTerm.trim();
      const data = await API.getProducts(params);
      const products = data.products || [];

      if (products.length === 0) {
        container.innerHTML = '<div style="padding:30px;text-align:center;color:var(--text-muted);">Không tìm thấy sản phẩm</div>';
        return;
      }

      const compareIds = new Set(this._items.map(p => p.product_id));
      const atMax = this._items.length >= this.MAX_COMPARE;

      container.innerHTML = products.map(p => {
        const img = (p.images && p.images[0]) || '';
        const added = compareIds.has(p.id);
        const disabled = atMax && !added;
        return `
          <div class="cms-product-card ${added ? 'added' : ''}">
            <img src="${img}" alt="${p.name}" />
            <div class="cms-pc-info">
              <div class="cms-pc-name">${p.name}</div>
              <div class="cms-pc-price">${formatPrice(p.price)}</div>
            </div>
            ${added
              ? '<span class="cms-pc-added"><i class="fa-solid fa-check"></i> Đã thêm</span>'
              : disabled
                ? '<span class="cms-pc-limit"><i class="fa-solid fa-ban"></i> Tối đa ' + this.MAX_COMPARE + '</span>'
                : `<button class="btn-primary" style="padding:6px 12px;font-size:12px;" onclick="Compare.add(${p.id}).then(() => { Compare._renderModalCompare(); Compare._renderProductSelector(); })">Thêm</button>`
            }
          </div>
        `;
      }).join('');
    } catch (err) {
      container.innerHTML = `<div style="padding:30px;text-align:center;color:var(--text-muted);">Lỗi: ${err.message}</div>`;
    }
  },

  // Search products in the selector
  async searchProducts(term) {
    await this._renderProductSelector(term);
  }
};
