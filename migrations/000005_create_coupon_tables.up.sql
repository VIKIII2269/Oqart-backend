-- Create ENUM types for coupons
CREATE TYPE coupon_type AS ENUM ('percentage', 'flat', 'free_shipping', 'buy_x_get_y', 'cashback');
CREATE TYPE coupon_status AS ENUM ('active', 'inactive', 'expired', 'deleted');
CREATE TYPE coupon_user_type AS ENUM ('all', 'new', 'existing', 'specific');

-- Coupons table
CREATE TABLE coupons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type coupon_type NOT NULL,
    status coupon_status NOT NULL DEFAULT 'active',
    user_type coupon_user_type NOT NULL DEFAULT 'all',

    -- Discount details
    discount_value DECIMAL(15, 2) NOT NULL,
    max_discount_amount DECIMAL(15, 2),
    min_order_value DECIMAL(15, 2) DEFAULT 0.00,

    -- Usage limits
    max_total_usage INT,
    current_usage INT DEFAULT 0,
    max_usage_per_user INT DEFAULT 1,

    -- Validity
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL,
    valid_until TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Conditions
    applicable_categories JSONB,
    applicable_products JSONB,
    applicable_vendors JSONB,
    excluded_categories JSONB,
    excluded_products JSONB,
    payment_methods JSONB,
    min_items INT DEFAULT 1,

    -- Buy X Get Y specific
    buy_quantity INT,
    get_quantity INT,
    buy_product_ids JSONB,
    get_product_ids JSONB,

    -- Metadata
    is_stackable BOOLEAN DEFAULT FALSE,
    is_auto_apply BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 0,

    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Coupon usage table
CREATE TABLE coupon_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coupon_id UUID NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID,
    discount_amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Coupon user restrictions table (for specific users)
CREATE TABLE coupon_user_restrictions (
    coupon_id UUID NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (coupon_id, user_id)
);

-- Indexes for coupons
CREATE INDEX idx_coupons_code ON coupons(code);
CREATE INDEX idx_coupons_status ON coupons(status);
CREATE INDEX idx_coupons_type ON coupons(type);
CREATE INDEX idx_coupons_valid_from ON coupons(valid_from);
CREATE INDEX idx_coupons_valid_until ON coupons(valid_until);
CREATE INDEX idx_coupons_is_public ON coupons(is_public);

-- Indexes for coupon_usage
CREATE INDEX idx_coupon_usage_coupon_id ON coupon_usage(coupon_id);
CREATE INDEX idx_coupon_usage_user_id ON coupon_usage(user_id);
CREATE INDEX idx_coupon_usage_order_id ON coupon_usage(order_id);

-- Indexes for coupon_user_restrictions
CREATE INDEX idx_coupon_user_restrictions_coupon_id ON coupon_user_restrictions(coupon_id);
CREATE INDEX idx_coupon_user_restrictions_user_id ON coupon_user_restrictions(user_id);

-- Add triggers for updated_at
CREATE TRIGGER update_coupons_updated_at BEFORE UPDATE ON coupons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to auto-update coupon status based on validity
CREATE OR REPLACE FUNCTION update_coupon_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.valid_until < CURRENT_TIMESTAMP AND NEW.status != 'expired' THEN
        NEW.status = 'expired';
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_coupon_status
    BEFORE INSERT OR UPDATE ON coupons
    FOR EACH ROW EXECUTE FUNCTION update_coupon_status();
