-- QMQSHOP Seed Data

-- ============================================================
-- Admin user (admin@qmqshop.com / 123456)
-- ============================================================
INSERT INTO users (email, password_hash, full_name, phone, address, role)
VALUES (
    'admin@qmqshop.com',
    '$2a$10$eKr6hfN/iNjJJmMDGzDEsOzklx8UNoT8chNoj7hwrXjg3m/upudpu',
    'Quản Trị Viên',
    '1900xxxx',
    'Hà Nội, Việt Nam',
    'admin'
);

-- ============================================================
-- Categories
-- ============================================================
INSERT INTO categories (name, slug, icon) VALUES
('Laptop Gaming',   'laptop-gaming',   'fa-laptop'),
('PC Gaming',       'pc-gaming',       'fa-gamepad'),
('PC Văn Phòng',    'pc-van-phong',    'fa-briefcase'),
('Linh Kiện',       'linh-kien',       'fa-microchip');

-- ============================================================
-- Products — Laptop Gaming (category_id = 1)
-- ============================================================
INSERT INTO products (category_id, name, slug, description, specs, price, old_price, images, stock, featured) VALUES
(1,
 'Laptop Gaming Acer Nitro V ANV15-51-57B4',
 'acer-nitro-v-anv15',
 'Acer Nitro V sở hữu hiệu năng mạnh mẽ với chip Intel Core i5-13420H, card đồ họa RTX 4050, RAM 16GB. Màn hình 15.6 inch Full HD 144Hz, thiết kế gaming đặc trưng với các đường cắt cơ khí. Hoàn hảo cho game thủ và người dùng đồ họa.',
 '{"CPU": "Intel Core i5-13420H", "RAM": "16GB DDR5", "GPU": "NVIDIA RTX 4050 6GB", "Màn hình": "15.6\" Full HD 144Hz", "SSD": "512GB NVMe", "Pin": "57Wh"}',
 21490000, 24990000,
 ARRAY['https://images.unsplash.com/photo-1603302576837-37561b2e2302?q=80&w=800&auto=format&fit=crop'],
 15, true),

(1,
 'Laptop ASUS ROG Strix G16 G614JU',
 'asus-rog-strix-g16',
 'ASUS ROG Strix G16 là cỗ máy gaming cao cấp với Intel Core i7-13650HX, RTX 4060. Màn hình 16 inch 2K 165Hz, hệ thống tản nhiệt ROG Intelligence. Đèn nền RGB đồng bộ AURA Sync. Lựa chọn tối thượng cho game thủ chuyên nghiệp.',
 '{"CPU": "Intel Core i7-13650HX", "RAM": "16GB DDR5", "GPU": "NVIDIA RTX 4060 8GB", "Màn hình": "16\" 2K 165Hz", "SSD": "1TB NVMe", "Pin": "90Wh"}',
 34990000, 38990000,
 ARRAY['https://images.unsplash.com/photo-1593642632823-8f785ba67e45?q=80&w=800&auto=format&fit=crop'],
 8, true),

(1,
 'Apple MacBook Pro 14 M3 Pro',
 'macbook-pro-14-m3',
 'MacBook Pro 14 inch với chip M3 Pro mang đến hiệu năng vượt trội cho công việc sáng tạo. Màn hình Liquid Retina XDR 14.2 inch, thời lượng pin lên đến 17 giờ. Thunderbolt 4, HDMI, SDXC. Thiết kế nhôm nguyên khối sang trọng.',
 '{"CPU": "Apple M3 Pro", "RAM": "18GB Unified", "GPU": "M3 Pro 18-core", "Màn hình": "14.2\" Liquid Retina XDR", "SSD": "512GB", "Pin": "17 giờ"}',
 49990000, 59990000,
 ARRAY['https://images.unsplash.com/photo-1541807084-5c52b6b3adef?q=80&w=800&auto=format&fit=crop'],
 5, true),

(1,
 'Laptop Dell XPS 15 9530',
 'dell-xps-15-9530',
 'Dell XPS 15 với thiết kế InfinityEdge màn hình gần như không viền. Chip Intel Core i7-13700H, RAM 16GB, RTX 4060. Màn hình OLED 3.5K cảm ứng. Thân máy bằng nhôm và sợi carbon siêu nhẹ. Cao cấp và tinh tế.',
 '{"CPU": "Intel Core i7-13700H", "RAM": "16GB DDR5", "GPU": "NVIDIA RTX 4060 8GB", "Màn hình": "15.6\" OLED 3.5K Touch", "SSD": "512GB NVMe", "Pin": "86Wh"}',
 39990000, 45990000,
 ARRAY['https://images.unsplash.com/photo-1593642702749-b7d2a804fbcf?q=80&w=800&auto=format&fit=crop'],
 3, true),

(1,
 'Laptop Lenovo Legion 5 Pro 16',
 'lenovo-legion-5-pro',
 'Lenovo Legion 5 Pro 16 inch với AMD Ryzen 7 7745HX, RTX 4070. Màn hình 16 inch WQXGA 240Hz, 100% sRGB. Hệ thống tản nhiệt ColdFront 5.0. Bàn phím TrueStrike với đèn RGB 4 vùng. Tối ưu cho gaming và stream.',
 '{"CPU": "AMD Ryzen 7 7745HX", "RAM": "32GB DDR5", "GPU": "NVIDIA RTX 4070 8GB", "Màn hình": "16\" WQXGA 240Hz", "SSD": "1TB NVMe", "Pin": "80Wh"}',
 42990000, 47990000,
 ARRAY['https://images.unsplash.com/photo-1629131726692-1accd0c53ce0?q=80&w=800&auto=format&fit=crop'],
 6, true);

-- ============================================================
-- Products — PC Gaming (category_id = 2)
-- ============================================================
INSERT INTO products (category_id, name, slug, description, specs, price, old_price, images, stock, featured) VALUES
(2,
 'PC Gaming MINH Vanguard i5/RTX 4060',
 'pc-minh-vanguard-i5',
 'PC Gaming MINH Vanguard được xây dựng trên nền tảng Intel Core i5-12400F và RTX 4060. Bo mạch chủ B760M, RAM 16GB DDR4. Case RGB kính cường lực. Hiệu năng chiến mọi tựa game eSports và AAA ở thiết lập High.',
 '{"CPU": "Intel Core i5-12400F", "RAM": "16GB DDR4 3200", "GPU": "NVIDIA RTX 4060 8GB", "Mainboard": "B760M DDR4", "SSD": "512GB NVMe", "PSU": "550W 80+"}',
 16500000, 18000000,
 ARRAY['https://images.unsplash.com/photo-1588872657578-7efd1f1555ed?q=80&w=800&auto=format&fit=crop'],
 20, true),

(2,
 'PC Gaming MINH Elite i7/RTX 4070',
 'pc-minh-elite-i7',
 'PC Gaming MINH Elite với Intel Core i7-13700KF và RTX 4070 12GB. RAM 32GB DDR5, SSD 1TB NVMe Gen4. Tản nhiệt AIO 240mm. PSU 750W Gold. Case cao cấp với 3 fan RGB. Chiến mọi tựa game 2K max setting.',
 '{"CPU": "Intel Core i7-13700KF", "RAM": "32GB DDR5 5600", "GPU": "NVIDIA RTX 4070 12GB", "Mainboard": "Z790 DDR5", "SSD": "1TB NVMe Gen4", "PSU": "750W 80+ Gold"}',
 28990000, 31990000,
 ARRAY['https://images.unsplash.com/photo-1587202372634-32705e3bf49c?q=80&w=800&auto=format&fit=crop'],
 10, true),

(2,
 'PC Gaming MINH Pro AMD/RTX 4080',
 'pc-minh-pro-amd',
 'PC Gaming MINH Pro đỉnh cao với AMD Ryzen 7 7800X3D và RTX 4080 Super 16GB. RAM 32GB DDR5 6000MHz, SSD 2TB NVMe Gen5. Tản nhiệt AIO 360mm. PSU 1000W Gold. Dàn máy chiến 4K Ultra mọi tựa game.',
 '{"CPU": "AMD Ryzen 7 7800X3D", "RAM": "32GB DDR5 6000", "GPU": "NVIDIA RTX 4080 Super 16GB", "Mainboard": "X670E DDR5", "SSD": "2TB NVMe Gen5", "PSU": "1000W 80+ Gold"}',
 45990000, 49990000,
 ARRAY['https://images.unsplash.com/photo-1587202372775-e229f172b9f7?q=80&w=800&auto=format&fit=crop'],
 4, false),

(2,
 'PC Gaming MINH Lite i3/RTX 3050',
 'pc-minh-lite-i3',
 'PC Gaming MINH Lite dành cho game thủ ngân sách hạn chế. Intel Core i3-12100F, RTX 3050 6GB, RAM 16GB DDR4. Chơi tốt Valorant, LOL, CS2, Fornite ở thiết lập High. Nâng cấp dễ dàng trong tương lai.',
 '{"CPU": "Intel Core i3-12100F", "RAM": "16GB DDR4 3200", "GPU": "NVIDIA RTX 3050 6GB", "Mainboard": "H610 DDR4", "SSD": "256GB NVMe", "PSU": "450W 80+"}',
 10990000, 12990000,
 ARRAY['https://images.unsplash.com/photo-1591488320449-011701bb6704?q=80&w=800&auto=format&fit=crop'],
 25, false);

-- ============================================================
-- Products — PC Văn Phòng (category_id = 3)
-- ============================================================
INSERT INTO products (category_id, name, slug, description, specs, price, old_price, images, stock, featured) VALUES
(3,
 'PC Văn Phòng MINH Office i5',
 'pc-office-i5',
 'PC Văn phòng MINH Office với Intel Core i5-12400, RAM 16GB, SSD 512GB. Card đồ họa tích hợp Intel UHD 730. Xử lý mượt mà Word, Excel, duyệt web, họp Zoom. Bảo hành 24 tháng. Hỗ trợ màn hình kép.',
 '{"CPU": "Intel Core i5-12400", "RAM": "16GB DDR4", "GPU": "Intel UHD 730", "Mainboard": "B660M", "SSD": "512GB NVMe", "PSU": "400W"}',
 8990000, 10500000,
 ARRAY['https://images.unsplash.com/photo-1593642632559-0c6d3fc62b89?q=80&w=800&auto=format&fit=crop'],
 30, true),

(3,
 'PC Văn Phòng MINH Pro i7',
 'pc-office-pro-i7',
 'PC Văn phòng MINH Pro với Intel Core i7-13700, RAM 32GB, SSD 1TB. Phù hợp cho dân thiết kế, kế toán, lập trình. Xử lý đa nhiệm cực tốt với 16 nhân 24 luồng. Mát mẻ, ổn định 24/7.',
 '{"CPU": "Intel Core i7-13700", "RAM": "32GB DDR5", "GPU": "Intel UHD 770", "Mainboard": "B760M", "SSD": "1TB NVMe", "PSU": "500W"}',
 14990000, 17500000,
 ARRAY['https://images.unsplash.com/photo-1593642634315-48f5414c3ad9?q=80&w=800&auto=format&fit=crop'],
 12, false),

(3,
 'PC Văn Phòng MINH Basic',
 'pc-office-basic',
 'PC Văn phòng MINH Basic - giải pháp giá rẻ cho văn phòng cơ bản. Intel Core i3-12100, RAM 8GB, SSD 256GB. Đáp ứng tốt nhu cầu làm việc văn phòng cơ bản, gửi mail, họp online. Tiết kiệm điện.',
 '{"CPU": "Intel Core i3-12100", "RAM": "8GB DDR4", "GPU": "Intel UHD 730", "Mainboard": "H610M", "SSD": "256GB NVMe", "PSU": "350W"}',
 5990000, 7500000,
 ARRAY['https://images.unsplash.com/photo-1593642634443-44adaa06623a?q=80&w=800&auto=format&fit=crop'],
 40, false);

-- ============================================================
-- Products — Linh Kiện (category_id = 4)
-- ============================================================
INSERT INTO products (category_id, name, slug, description, specs, price, old_price, images, stock, featured) VALUES
(4,
 'Card Đồ Họa NVIDIA RTX 4070 Super',
 'nvidia-rtx-4070-super',
 'NVIDIA GeForce RTX 4070 SUPER với 12GB GDDR6X. Kiến trúc Ada Lovelace, DLSS 3.5, Ray Tracing thế hệ 3. Hiệu năng vượt trội so với thế hệ trước. Phù hợp gaming 2K và 4K, dựng hình, render AI.',
 '{"VRAM": "12GB GDDR6X", "CUDA Cores": "7168", "Bus": "192-bit", "TDP": "220W", "Kết nối": "HDMI 2.1 + 3x DisplayPort 1.4a", "Kích thước": "2.5 slot"}',
 14990000, 16990000,
 ARRAY['https://images.unsplash.com/photo-1591488320449-011701bb6704?q=80&w=800&auto=format&fit=crop'],
 7, true),

(4,
 'CPU Intel Core i5-14600K',
 'intel-i5-14600k',
 'Intel Core i5-14600K thế hệ 14 với 14 nhân (6 P-core + 8 E-core), 20 luồng. Turbo Boost tối đa 5.3GHz. Hỗ trợ DDR5 và DDR4. Socket LGA 1700. Tương thích mainboard B760, Z790. Hiệu năng gaming xuất sắc.',
 '{"Nhân/Luồng": "14/20 (6P+8E)", "Xung tối đa": "5.3 GHz", "L3 Cache": "24MB", "TDP": "125W (181W Turbo)", "Socket": "LGA 1700", "Hỗ trợ RAM": "DDR4/DDR5"}',
 7290000, 8290000,
 ARRAY['https://images.unsplash.com/photo-1591799264318-7e6ef8ddb7ea?q=80&w=800&auto=format&fit=crop'],
 18, true),

(4,
 'RAM Kingston Fury 16GB DDR5 5600',
 'kingston-fury-16gb-ddr5',
 'RAM Kingston Fury Beast 16GB DDR5 5600MHz. Tản nhiệt nhôm cao cấp, hỗ trợ Intel XMP 3.0. Plug-and-play, tương thích đa nền tảng. Nâng cấp hiệu năng đáng kể so với DDR4. Bảo hành trọn đời.',
 '{"Dung lượng": "16GB (1x16GB)", "Loại": "DDR5", "Tốc độ": "5600MHz", "Điện áp": "1.25V", "CAS Latency": "CL40", "Tản nhiệt": "Nhôm ép"}',
 1190000, 1490000,
 ARRAY['https://images.unsplash.com/photo-1606914504570-53533a4e1bd3?q=80&w=800&auto=format&fit=crop'],
 50, false),

(4,
 'SSD Samsung 990 Pro 1TB NVMe Gen4',
 'samsung-990-pro-1tb',
 'SSD Samsung 990 PRO 1TB PCIe 4.0 NVMe M.2 với tốc độ đọc lên đến 7,450MB/s. Tối ưu cho gaming và xử lý đồ họa nặng. Bộ điều khiển Samsung Pascal tự động. Bảo hành 5 năm.',
 '{"Dung lượng": "1TB", "Chuẩn": "PCIe 4.0 x4 NVMe M.2", "Đọc": "7,450 MB/s", "Ghi": "6,900 MB/s", "TBW": "600TB", "Bảo hành": "5 năm"}',
 3490000, 4290000,
 ARRAY['https://images.unsplash.com/photo-1563770551464-0d9d4e8ad380?q=80&w=800&auto=format&fit=crop'],
 22, false);
