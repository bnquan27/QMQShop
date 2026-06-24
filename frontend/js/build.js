const COMPONENT_LABELS = {
	cpu: "CPU",
	mainboard: "Mainboard",
	ram: "RAM",
	gpu: "Card Đồ Họa",
	storage: "Ổ Cứng",
	cooler: "Tản Nhiệt",
	case: "Vỏ Case",
	monitor: "Màn Hình",
};

const COMPONENT_ORDER = [
	"cpu",
	"mainboard",
	"ram",
	"gpu",
	"storage",
	"cooler",
	"case",
	"monitor",
];

let buildData = {};
let allComponents = {};

function loadBuild() {
	try {
		const saved = localStorage.getItem("qmq_build");
		if (saved) buildData = JSON.parse(saved);
		else buildData = {};
	} catch {
		buildData = {};
	}
}

function saveBuild() {
	localStorage.setItem("qmq_build", JSON.stringify(buildData));
}

function showDangerPopup(title, message, buttons) {
	document.getElementById("danger-popup-title").textContent = title;
	document.getElementById("danger-popup-message").textContent = message;
	document.getElementById("danger-popup-actions").innerHTML = buttons
		.map(function(b) {
			return (
				'<button class="' +
				(b.btnClass || "btn-outline") +
				' btn-sm" onclick="' +
				b.onClick +
				'">' +
				b.label +
				"</button>"
			);
		})
		.join("");
	document.getElementById("danger-popup").style.display = "flex";
	document.body.style.overflow = "hidden";
}

function closeDangerPopup() {
	document.getElementById("danger-popup").style.display = "none";
	document.body.style.overflow = "";
}

function clearBuild() {
	showDangerPopup(
		"Xóa cấu hình",
		"Cảnh báo: Toàn bộ linh kiện của bộ PC hiện tại sẽ bị xóa đi",
		[
			{ label: "Hủy", onClick: "closeDangerPopup()" },
			{
				label: "Xác nhận",
				btnClass: "btn-danger",
				onClick: "confirmClearBuild()",
			},
		],
	);
}

function confirmClearBuild() {
	closeDangerPopup();
	buildData = {};
	saveBuild();
	renderAll();
}

function addBuildToCartWithCheck() {
	var count = getSelectedCount();
	if (count === 0) {
		showDangerPopup("Oops...", "Bạn chưa chọn sản phẩm nào", [
			{ label: "Đóng", onClick: "closeDangerPopup()" },
		]);
		return;
	}
	addBuildToCart();
}

function esc(s) {
	return String(s)
		.replace(/&/g, "&amp;")
		.replace(/</g, "&lt;")
		.replace(/>/g, "&gt;")
		.replace(/"/g, "&quot;")
		.replace(/'/g, "&#39;");
}

function formatPrice(v) {
	return (v || 0).toLocaleString("vi-VN") + "₫";
}

function getTypeTotal(type) {
	const p = buildData[type];
	return p ? p.price * (p.qty || 1) : 0;
}

function getTotalPrice() {
	return COMPONENT_ORDER.reduce((sum, t) => sum + getTypeTotal(t), 0);
}

function getSelectedCount() {
	return COMPONENT_ORDER.filter((t) => buildData[t]).length;
}

function selectComponent(type, productId) {
	const products = allComponents[type] || [];
	const product = products.find((p) => p.id === productId);
	if (!product) return;
	buildData[type] = Object.assign({}, product, { qty: 1 });
	saveBuild();
	closeComponentPopup();
	renderAll();
}

function removeComponent(type) {
	delete buildData[type];
	saveBuild();
	renderAll();
}

function changeQty(type, delta) {
	const p = buildData[type];
	if (!p) return;
	const qty = Math.max(1, (p.qty || 1) + delta);
	p.qty = qty;
	saveBuild();
	renderAll();
}

function renderSummary() {
	const count = getSelectedCount();
	document.getElementById("summary-count").textContent = count + "/8 linh kiện";
	document.getElementById("summary-total").textContent =
		"Chi phí ước tính: " + formatPrice(getTotalPrice());
	document.getElementById("summary-add-cart").style.display =
		count > 0 ? "" : "none";
}

async function addBuildToCart() {
	const selected = COMPONENT_ORDER.filter((t) => buildData[t]);
	if (selected.length === 0) return;
	let added = 0;
	for (const t of selected) {
		const item = buildData[t];
		try {
			await API.request("POST", "/api/cart", {
				product_id: item.id,
				quantity: item.qty || 1,
			});
			added += item.qty || 1;
		} catch (err) {
			console.error("Add to cart failed for", item.name, err);
		}
	}
	if (added > 0) {
		await Cart.load();
		window.location.href = "/cart.html";
	} else {
		alert("Không thể thêm vào giỏ hàng. Vui lòng đăng nhập.");
	}
}

function scrollToFirstEmpty() {
	for (const t of COMPONENT_ORDER) {
		if (!buildData[t]) {
			const el = document.getElementById("section-" + t);
			if (el) el.scrollIntoView({ behavior: "smooth", block: "start" });
			return;
		}
	}
}

function renderAll() {
	const container = document.getElementById("component-sections");
	let html = "";

	for (const t of COMPONENT_ORDER) {
		const label = COMPONENT_LABELS[t] || t;
		const selected = buildData[t];

		html +=
			'<div class="builder-section" id="section-' +
			t +
			'">' +
			'<div class="builder-section-header">' +
			'<div class="builder-section-title">' +
			'<span class="builder-section-badge">' +
			label +
			"</span>" +
			"</div>" +
			'<div class="builder-section-actions">' +
			(selected
				? '<button class="btn-primary btn-sm" onclick="openComponentPopup(\'' +
				t +
				'\')"><i class="fa-solid fa-arrow-right-arrow-left"></i> Đổi</button>'
				: '<button class="btn-primary btn-sm" onclick="openComponentPopup(\'' +
				t +
				'\')"><i class="fa-solid fa-plus"></i> Chọn</button>') +
			"</div>" +
			"</div>" +
			(selected
				? '<div class="builder-section-body">' +
				'<div class="item-left">' +
				'<img class="item-left-img" src="' +
				esc((selected.images && selected.images[0]) || "") +
				'" alt="' +
				esc(selected.name) +
				'" loading="lazy" />' +
				'<div class="item-left-info">' +
				'<a class="item-left-name" href="/product.html?id=' +
				selected.id +
				'">' +
				esc(selected.name) +
				"</a>" +
				'<div class="item-left-status ' +
				(selected.stock > 0 ? "in-stock" : "out-of-stock") +
				'">' +
				(selected.stock > 0 ? "Còn hàng" : "Hết hàng") +
				"</div>" +
				"</div>" +
				"</div>" +
				'<div class="item-right">' +
				'<div class="item-right-row"><span class="item-right-value">' +
				formatPrice(selected.price) +
				"</span></div>" +
				"<span style='margin: 10px'>x</span>" +
				'<div class="item-right-row">' +
				'<span class="item-qty-spinner">' +
				'<button class="qty-btn qty-minus" onclick="changeQty(\'' +
				t +
				'\', -1); return false;"><i class=\"fa-solid fa-minus\"></i></button>' +
				'<span class="qty-value">' +
				(selected.qty || 1) +
				"</span>" +
				'<button class="qty-btn qty-plus" onclick="changeQty(\'' +
				t +
				'\', 1); return false;"><i class=\"fa-solid fa-plus\"></i></button>' +
				"</span></div>" +
				'<div class="item-right-row item-right-total-row"><span>=</span><span class="item-right-total">' +
				formatPrice((selected.price || 0) * (selected.qty || 1)) +
				"</span></div>" +
				'<button class="btn-ghost btn-sm item-right-remove" onclick="removeComponent(\'' +
				t +
				'\')\" title=\"Bỏ chọn\"><i class=\"fa-regular fa-trash-can\"></i> Xóa</button>' +
				"</div>" +
				"</div>"
				: '<div class="builder-section-body builder-section-body-empty"><span class="builder-section-empty">Chưa chọn</span></div>') +
			"</div>";
	}

	container.innerHTML = html;
	renderSummary();
}

// --- Popup --------------------------------------------------
// Build filter definitions dynamically from product specs
function getFiltersForType() {
	const products = popupState.products || [];
	const specKeys = [];
	const seen = {};
	for (const p of products) {
		if (p.specs) {
			for (const key of Object.keys(p.specs)) {
				if (!seen[key]) {
					seen[key] = true;
					specKeys.push(key);
				}
			}
		}
	}
	specKeys.sort();

	const filters = [
		{ id: "brand", label: "Hãng sản xuất", field: "brand" },
	];
	for (const key of specKeys) {
		filters.push({
			id: "specs_" + key.replace(/[^a-zA-Z0-9_À-ỹ]/g, "_").replace(/_+/g, "_").replace(/^_|_$/g, ""),
			label: key,
			field: "specs->" + key,
		});
	}
	return filters;
}

let popupState = {
	type: null,
	products: [],
	filtered: [],
	page: 1,
	perPage: 6,
	activeFilters: {},
};

function getFilterValue(product, filterDef) {
	if (filterDef.field === "brand") return product.brand || "";
	if (filterDef.field.startsWith("specs->")) {
		const key = filterDef.field.substring(7);
		return (product.specs && product.specs[key]) || "";
	}
	return "";
}

function openComponentPopup(type) {
	const products = allComponents[type] || [];
	if (products.length === 0) return;

	popupState.type = type;
	popupState.products = products;
	popupState.page = 1;
	popupState.activeFilters = {};

	const label = COMPONENT_LABELS[type] || type;
	document.getElementById("popup-title").textContent = "Chọn " + label;
	document.getElementById("popup-sort").value = "newest";
	document.getElementById("popup-stock-only").checked = false;

	renderPopupFilters();
	applyPopupFilters();

	document.getElementById("component-popup").style.display = "flex";
	document.body.style.overflow = "hidden";
}

function closeComponentPopup() {
	document.getElementById("component-popup").style.display = "none";
	document.body.style.overflow = "";
}

function renderPopupFilters() {
	const type = popupState.type;
	const products = popupState.products;
	const filters = getFiltersForType();
	const container = document.getElementById("popup-filters");
	const activeFilters = popupState.activeFilters;

	let html = '<div class="popup-filters-inner">';
	for (const f of filters) {
		const values = [
			...new Set(
				products
					.map(function(p) {
						return getFilterValue(p, f);
					})
					.filter(function(v) {
						return v;
					}),
			),
		].sort();
		if (values.length === 0) continue;

		html +=
			'<div class="popup-filter-group"><div class="popup-filter-label">' +
			esc(f.label) +
			"</div>";
		for (const v of values) {
			const checked =
				activeFilters[f.id] && activeFilters[f.id].indexOf(v) !== -1;
			html +=
				'<label class="popup-filter-option"><input type="checkbox" data-filter-id="' +
				esc(f.id) +
				'" value="' +
				esc(v) +
				'"' +
				(checked ? " checked" : "") +
				" onchange=\"togglePopupFilter('" +
				esc(f.id) +
				"','" +
				esc(v) +
				"')\" /><span>" +
				esc(v) +
				"</span></label>";
		}
		html += "</div>";
	}
	html += "</div>";
	container.innerHTML =
		html || '<div class="popup-filter-empty">Không có bộ lọc</div>';
}

function togglePopupFilter(filterId, value) {
	const active = popupState.activeFilters;
	if (!active[filterId]) active[filterId] = [];
	const idx = active[filterId].indexOf(value);
	if (idx >= 0) active[filterId].splice(idx, 1);
	else active[filterId].push(value);
	if (active[filterId].length === 0) delete active[filterId];
	applyPopupFilters();
}

function applyPopupFilters() {
	const type = popupState.type;
	const products = popupState.products;
	const sort = document.getElementById("popup-sort").value;
	const stockOnly = document.getElementById("popup-stock-only").checked;
	const activeFilters = popupState.activeFilters;
	const filters = getFiltersForType();

	let filtered = products.slice();
	if (stockOnly)
		filtered = filtered.filter(function(p) {
			return p.stock > 0;
		});

	for (const filterId of Object.keys(activeFilters)) {
		const vals = activeFilters[filterId];
		if (!vals || vals.length === 0) continue;
		let filterDef = null;
		for (const f of filters) {
			if (f.id === filterId) {
				filterDef = f;
				break;
			}
		}
		if (!filterDef) continue;
		filtered = filtered.filter(function(p) {
			return vals.indexOf(getFilterValue(p, filterDef)) !== -1;
		});
	}

	if (sort === "newest")
		filtered.sort(function(a, b) {
			return new Date(b.created_at) - new Date(a.created_at);
		});
	else if (sort === "price-asc")
		filtered.sort(function(a, b) {
			return a.price - b.price;
		});
	else if (sort === "price-desc")
		filtered.sort(function(a, b) {
			return b.price - a.price;
		});

	popupState.filtered = filtered;
	popupState.page = 1;
	renderPopupProducts();
}

function clearPopupFilters() {
	popupState.activeFilters = {};
	document.getElementById("popup-sort").value = "newest";
	renderPopupFilters();
	applyPopupFilters();
}

function renderPopupProducts() {
	const filtered = popupState.filtered;
	const page = popupState.page;
	const perPage = popupState.perPage;
	const totalPages = Math.ceil(filtered.length / perPage) || 1;
	const start = (page - 1) * perPage;
	const pageProducts = filtered.slice(start, start + perPage);
	const selected = buildData[popupState.type];

	let pagHtml = "";
	if (totalPages > 1) {
		pagHtml +=
			'<button class="btn-ghost btn-sm pag-btn" onclick="changePopupPage(-1)"' +
			(page <= 1 ? " disabled" : "") +
			'><i class="fa-solid fa-chevron-left"></i></button>';
		for (let i = 1; i <= totalPages; i++) {
			pagHtml +=
				'<button class="btn-ghost btn-sm pag-btn' +
				(i === page ? " active" : "") +
				'" onclick="goToPopupPage(' +
				i +
				')">' +
				i +
				"</button>";
		}
		pagHtml +=
			'<button class="btn-ghost btn-sm pag-btn" onclick="changePopupPage(1)"' +
			(page >= totalPages ? " disabled" : "") +
			'><i class="fa-solid fa-chevron-right"></i></button>';
	}
	document.getElementById("popup-pagination").innerHTML = pagHtml;

	const container = document.getElementById("popup-products");
	if (pageProducts.length === 0) {
		container.innerHTML =
			'<div class="popup-empty">Không có sản phẩm phù hợp với bộ lọc.</div>';
		return;
	}

	let html = '<div class="popup-product-list">';
	for (const p of pageProducts) {
		const isSelected = selected && selected.id === p.id;
		const inStock = p.stock > 0;
		const img = (p.images && p.images[0]) || "";
		html +=
			'<div class="popup-product-item' +
			(isSelected ? " selected" : "") +
			'">' +
			'<div class="popup-product-img"><img src="' +
			esc(img) +
			'" alt="' +
			esc(p.name) +
			'" loading="lazy" /></div>' +
			'<div class="popup-product-info">' +
			'<div class="popup-product-name">' +
			esc(p.name) +
			"</div>" +
			'<div class="popup-product-status ' +
			(inStock ? "in-stock" : "out-of-stock") +
			'">' +
			(inStock ? "Còn hàng" : "Hết hàng") +
			"</div>" +
			'<div class="popup-product-price">' +
			formatPrice(p.price) +
			"</div>" +
			"</div>" +
			'<button class="btn-primary btn-sm popup-select-btn" onclick="selectComponent(\'' +
			popupState.type +
			"', " +
			p.id +
			')"' +
			(isSelected ? " disabled" : "") +
			">" +
			(isSelected ? "Đã chọn" : "Thêm vào cấu hình") +
			"</button>" +
			"</div>";
	}
	html += "</div>";
	container.innerHTML = html;
}

function changePopupPage(delta) {
	const newPage = popupState.page + delta;
	const totalPages = Math.ceil(popupState.filtered.length / popupState.perPage);
	if (newPage < 1 || newPage > totalPages) return;
	popupState.page = newPage;
	renderPopupProducts();
}

function goToPopupPage(page) {
	popupState.page = page;
	renderPopupProducts();
}

async function loadComponents() {
	try {
		const groups = await API.request("GET", "/api/pc-builder/components");
		allComponents = {};
		for (const g of groups) {
			allComponents[g.type] = g.products;
		}
		renderAll();
	} catch (err) {
		document.getElementById("component-sections").innerHTML =
			'<div class="loading-center" style="grid-column:1/-1;">' +
			'<i class="fa-regular fa-face-frown" style="font-size:40px;color:var(--text-muted);"></i>' +
			'<span style="font-size:16px;color:var(--text-secondary);">Lỗi tải linh kiện: ' +
			esc(err.message) +
			"</span>" +
			"</div>";
	}
}

document.addEventListener("DOMContentLoaded", async () => {
	await Auth.init();
	Auth.updateHeader();
	await Cart.load();
	loadBuild();
	await loadComponents();
});
