-- Create ENUM types for orders
CREATE TYPE order_status AS ENUM ('pending', 'payment_confirmed', 'accepted', 'rejected', 'processing', 'packed', 'shipped', 'out_for_delivery', 'delivered', 'cancelled', 'refunded');
CREATE TYPE payment_method AS ENUM ('cod', 'razorpay', 'stripe', 'wallet', 'upi', 'card', 'netbanking');
CREATE TYPE payment_status AS ENUM ('pending', 'processing', 'completed', 'failed', 'refunded', 'partially_refunded');

-- Orders table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE RESTRICT,
    status order_status NOT NULL DEFAULT 'pending',

    -- Amounts
    subtotal DECIMAL(15, 2) NOT NULL,
    tax DECIMAL(15, 2) DEFAULT 0.00,
    shipping_charge DECIMAL(15, 2) DEFAULT 0.00,
    coupon_discount DECIMAL(15, 2) DEFAULT 0.00,
    total DECIMAL(15, 2) NOT NULL,

    -- Coupon
    coupon_id UUID REFERENCES coupons(id),
    coupon_code VARCHAR(50),

    -- Address
    shipping_address JSONB NOT NULL,
    billing_address JSONB NOT NULL,

    -- Payment
    payment_method payment_method NOT NULL,
    payment_status payment_status NOT NULL DEFAULT 'pending',

    -- Delivery
    estimated_delivery_date DATE,
    actual_delivery_date TIMESTAMP WITH TIME ZONE,
    tracking_number VARCHAR(255),
    tracking_url VARCHAR(500),
    delivery_instructions TEXT,

    -- Customer details
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(20) NOT NULL,

    -- Additional info
    notes TEXT,
    cancellation_reason TEXT,
    rejection_reason TEXT,
    cancelled_by UUID REFERENCES users(id),
    cancelled_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Order items table
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    variant_id UUID REFERENCES product_variants(id) ON DELETE RESTRICT,

    -- Product details snapshot
    product_name VARCHAR(255) NOT NULL,
    product_sku VARCHAR(100) NOT NULL,
    variant_name VARCHAR(255),
    product_image VARCHAR(500),

    -- Pricing
    quantity INT NOT NULL CHECK (quantity > 0),
    price DECIMAL(15, 2) NOT NULL,
    mrp DECIMAL(15, 2) NOT NULL,
    discount_percentage DECIMAL(5, 2) DEFAULT 0.00,
    tax_percentage DECIMAL(5, 2) DEFAULT 0.00,
    subtotal DECIMAL(15, 2) NOT NULL,

    -- Additional details
    weight DECIMAL(10, 2),
    weight_unit VARCHAR(20),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Order status history table
CREATE TABLE order_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status order_status NOT NULL,
    notes TEXT,
    changed_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    transaction_id VARCHAR(255) UNIQUE,
    gateway_order_id VARCHAR(255),
    payment_method payment_method NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',

    -- Amounts
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'INR',

    -- Gateway specific
    gateway VARCHAR(50),
    gateway_response JSONB,
    signature VARCHAR(500),

    -- Metadata
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    error_code VARCHAR(100),
    error_message TEXT,

    -- Refund details
    refund_amount DECIMAL(15, 2) DEFAULT 0.00,
    refund_status VARCHAR(50),
    refund_reference VARCHAR(255),
    refunded_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Payment transactions table (for tracking all payment attempts)
CREATE TABLE payment_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    transaction_id VARCHAR(255),
    status payment_status NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    gateway_response JSONB,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Invoices table
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    invoice_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE,
    subtotal DECIMAL(15, 2) NOT NULL,
    tax DECIMAL(15, 2) DEFAULT 0.00,
    shipping_charge DECIMAL(15, 2) DEFAULT 0.00,
    discount DECIMAL(15, 2) DEFAULT 0.00,
    total DECIMAL(15, 2) NOT NULL,
    pdf_path VARCHAR(500),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for orders
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_vendor_id ON orders(vendor_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_user_status ON orders(user_id, status);
CREATE INDEX idx_orders_vendor_status ON orders(vendor_id, status);

-- Indexes for order_items
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- Indexes for order_status_history
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_order_status_history_status ON order_status_history(status);

-- Indexes for payments
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_gateway_order_id ON payments(gateway_order_id);

-- Indexes for payment_transactions
CREATE INDEX idx_payment_transactions_payment_id ON payment_transactions(payment_id);

-- Indexes for invoices
CREATE INDEX idx_invoices_order_id ON invoices(order_id);
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);

-- Add triggers for updated_at
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_invoices_updated_at BEFORE UPDATE ON invoices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to generate order number
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.order_number IS NULL THEN
        NEW.order_number = 'ORD' || TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDD') || LPAD(nextval('order_number_seq')::TEXT, 6, '0');
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create sequence for order numbers
CREATE SEQUENCE order_number_seq START 1;

CREATE TRIGGER trigger_generate_order_number
    BEFORE INSERT ON orders
    FOR EACH ROW EXECUTE FUNCTION generate_order_number();

-- Function to generate invoice number
CREATE OR REPLACE FUNCTION generate_invoice_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.invoice_number IS NULL THEN
        NEW.invoice_number = 'INV' || TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDD') || LPAD(nextval('invoice_number_seq')::TEXT, 6, '0');
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create sequence for invoice numbers
CREATE SEQUENCE invoice_number_seq START 1;

CREATE TRIGGER trigger_generate_invoice_number
    BEFORE INSERT ON invoices
    FOR EACH ROW EXECUTE FUNCTION generate_invoice_number();

-- Function to automatically create order status history
CREATE OR REPLACE FUNCTION create_order_status_history()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' AND OLD.status != NEW.status) THEN
        INSERT INTO order_status_history (order_id, status, notes)
        VALUES (NEW.id, NEW.status, NULL);
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_create_order_status_history
    AFTER INSERT OR UPDATE OF status ON orders
    FOR EACH ROW EXECUTE FUNCTION create_order_status_history();

-- Function to update vendor order count and sales
CREATE OR REPLACE FUNCTION update_vendor_order_stats()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'delivered' AND (TG_OP = 'INSERT' OR OLD.status != 'delivered') THEN
        UPDATE vendors
        SET total_orders = total_orders + 1,
            total_sales = total_sales + NEW.total
        WHERE id = NEW.vendor_id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_vendor_order_stats
    AFTER INSERT OR UPDATE OF status ON orders
    FOR EACH ROW EXECUTE FUNCTION update_vendor_order_stats();
