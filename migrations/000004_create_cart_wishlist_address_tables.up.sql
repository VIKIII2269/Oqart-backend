-- Create ENUM types
CREATE TYPE address_type AS ENUM ('home', 'work', 'other');

-- Addresses table
CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type address_type NOT NULL DEFAULT 'home',
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    street_address VARCHAR(255) NOT NULL,
    landmark VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    pincode VARCHAR(10) NOT NULL,
    country VARCHAR(100) DEFAULT 'India',
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Carts table
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(255),
    coupon_code VARCHAR(50),
    coupon_discount DECIMAL(15, 2) DEFAULT 0.00,
    subtotal DECIMAL(15, 2) DEFAULT 0.00,
    tax DECIMAL(15, 2) DEFAULT 0.00,
    shipping_charge DECIMAL(15, 2) DEFAULT 0.00,
    total DECIMAL(15, 2) DEFAULT 0.00,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id),
    UNIQUE(session_id)
);

-- Cart items table
CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cart_id UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    price DECIMAL(15, 2) NOT NULL,
    mrp DECIMAL(15, 2) NOT NULL,
    discount_percentage DECIMAL(5, 2) DEFAULT 0.00,
    subtotal DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cart_id, product_id, variant_id)
);

-- Wishlists table
CREATE TABLE wishlists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Wishlist items table
CREATE TABLE wishlist_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wishlist_id UUID NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(wishlist_id, product_id)
);

-- Indexes for addresses
CREATE INDEX idx_addresses_user_id ON addresses(user_id);
CREATE INDEX idx_addresses_pincode ON addresses(pincode);
CREATE INDEX idx_addresses_is_default ON addresses(is_default);

-- Indexes for carts
CREATE INDEX idx_carts_user_id ON carts(user_id);
CREATE INDEX idx_carts_session_id ON carts(session_id);
CREATE INDEX idx_carts_expires_at ON carts(expires_at);

-- Indexes for cart_items
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);
CREATE INDEX idx_cart_items_variant_id ON cart_items(variant_id);

-- Indexes for wishlists
CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);

-- Indexes for wishlist_items
CREATE INDEX idx_wishlist_items_wishlist_id ON wishlist_items(wishlist_id);
CREATE INDEX idx_wishlist_items_product_id ON wishlist_items(product_id);

-- Add triggers for updated_at
CREATE TRIGGER update_addresses_updated_at BEFORE UPDATE ON addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_carts_updated_at BEFORE UPDATE ON carts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cart_items_updated_at BEFORE UPDATE ON cart_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_wishlists_updated_at BEFORE UPDATE ON wishlists
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to calculate cart item subtotal
CREATE OR REPLACE FUNCTION calculate_cart_item_subtotal()
RETURNS TRIGGER AS $$
BEGIN
    NEW.subtotal = NEW.price * NEW.quantity;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_calculate_cart_item_subtotal
    BEFORE INSERT OR UPDATE OF quantity, price ON cart_items
    FOR EACH ROW EXECUTE FUNCTION calculate_cart_item_subtotal();

-- Function to update cart totals
CREATE OR REPLACE FUNCTION update_cart_totals()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE carts
    SET subtotal = (
        SELECT COALESCE(SUM(subtotal), 0)
        FROM cart_items
        WHERE cart_id = NEW.cart_id
    ),
    total = (
        SELECT COALESCE(SUM(subtotal), 0) + COALESCE(carts.tax, 0) + COALESCE(carts.shipping_charge, 0) - COALESCE(carts.coupon_discount, 0)
        FROM cart_items
        WHERE cart_id = NEW.cart_id
    )
    WHERE id = NEW.cart_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_cart_totals_on_item_change
    AFTER INSERT OR UPDATE OR DELETE ON cart_items
    FOR EACH ROW EXECUTE FUNCTION update_cart_totals();

-- Function to ensure only one default address per user
CREATE OR REPLACE FUNCTION ensure_one_default_address()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_default = TRUE THEN
        UPDATE addresses
        SET is_default = FALSE
        WHERE user_id = NEW.user_id AND id != NEW.id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_ensure_one_default_address
    BEFORE INSERT OR UPDATE OF is_default ON addresses
    FOR EACH ROW EXECUTE FUNCTION ensure_one_default_address();
