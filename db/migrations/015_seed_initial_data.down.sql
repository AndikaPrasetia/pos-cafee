-- Delete seeded data (in reverse order due to foreign keys)
DELETE FROM inventory;
DELETE FROM menu_items;
DELETE FROM users;
DELETE FROM categories;
