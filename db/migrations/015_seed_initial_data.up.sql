-- Insert default categories
INSERT INTO categories (name, description) VALUES
('Coffee', 'Various coffee drinks'),
('Non-Coffee', 'Non-caffeinated beverages'),
('Snacks', 'Light snacks and pastries'),
('Meals', 'Heavy meals and lunch options');

-- Insert default users (passwords are hashed 'password123')
INSERT INTO users (username, password_hash, role, full_name) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMye.vb.8VYV2bZ.KDnR1J1vjYVqQ38ZCyW', 'admin', 'System Administrator'),
('manager', '$2a$10$N9qo8uLOickgx2ZMRZoMye.vb.8VYV2bZ.KDnR1J1vjYVqQ38ZCyW', 'manager', 'Cafe Manager'),
('kasir1', '$2a$10$N9qo8uLOickgx2ZMRZoMye.vb.8VYV2bZ.KDnR1J1vjYVqQ38ZCyW', 'kasir', 'Cashier One');

-- Insert sample menu items
INSERT INTO menu_items (name, category_id, price, cost, description) VALUES
('Espresso', 1, 15000, 5000, 'Strong black coffee'),
('Cappuccino', 1, 25000, 8000, 'Coffee with milk foam'),
('Latte', 1, 27000, 9000, 'Coffee with steamed milk'),
('Milk Tea', 2, 20000, 6000, 'Sweet milk tea'),
('Croissant', 3, 18000, 7000, 'Buttery French croissant'),
('Sandwich', 4, 35000, 15000, 'Club sandwich with chicken');

-- Insert initial inventory
INSERT INTO inventory (menu_item_id, current_stock, min_stock) VALUES
(1, 100, 20), (2, 80, 15), (3, 75, 15),
(4, 60, 10), (5, 40, 5), (6, 30, 5);
