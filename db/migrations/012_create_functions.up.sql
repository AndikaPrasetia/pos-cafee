-- Function to update timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to generate order number
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TRIGGER AS $$
DECLARE
    today_date VARCHAR(8);
    sequence_num INTEGER;
    new_order_number VARCHAR(20);
BEGIN
    today_date := TO_CHAR(CURRENT_DATE, 'YYYYMMDD');

    SELECT COALESCE(MAX(CAST(SUBSTRING(order_number FROM 10) AS INTEGER)), 0) + 1
    INTO sequence_num
    FROM orders
    WHERE order_number LIKE 'ORD-' || today_date || '-%';

    new_order_number := 'ORD-' || today_date || '-' || LPAD(sequence_num::TEXT, 4, '0');
    NEW.order_number := new_order_number;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to handle inventory updates on order completion
CREATE OR REPLACE FUNCTION update_inventory_on_order()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'completed' AND OLD.status != 'completed' THEN
        UPDATE inventory
        SET current_stock = current_stock - oi.quantity,
            last_updated = CURRENT_TIMESTAMP,
            updated_by = NEW.user_id
        FROM order_items oi
        WHERE inventory.menu_item_id = oi.menu_item_id
        AND oi.order_id = NEW.id;

        INSERT INTO stock_transactions (
            inventory_id, transaction_type, quantity, previous_stock,
            new_stock, reason, reference_id, reference_type, user_id
        )
        SELECT
            inv.id,
            'out',
            oi.quantity,
            inv.current_stock,
            (inv.current_stock - oi.quantity),
            'order_usage',
            NEW.id,
            'order',
            NEW.user_id
        FROM order_items oi
        JOIN inventory inv ON oi.menu_item_id = inv.menu_item_id
        WHERE oi.order_id = NEW.id;

        NEW.completed_at := CURRENT_TIMESTAMP;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to calculate order total
CREATE OR REPLACE FUNCTION calculate_order_total()
RETURNS TRIGGER AS $$
BEGIN
    NEW.subtotal := NEW.quantity * NEW.unit_price;

    UPDATE orders
    SET total_amount = (
        SELECT COALESCE(SUM(subtotal), 0)
        FROM order_items
        WHERE order_id = NEW.order_id
    ),
    updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.order_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to validate stock availability
CREATE OR REPLACE FUNCTION validate_stock_availability()
RETURNS TRIGGER AS $$
DECLARE
    current_stock_val INTEGER;
BEGIN
    IF NEW.status = 'completed' AND OLD.status != 'completed' THEN
        FOR current_stock_val IN
            SELECT inv.current_stock
            FROM order_items oi
            JOIN inventory inv ON oi.menu_item_id = inv.menu_item_id
            WHERE oi.order_id = NEW.id
        LOOP
            IF current_stock_val < 0 THEN
                RAISE EXCEPTION 'Insufficient stock for order completion';
            END IF;
        END LOOP;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
