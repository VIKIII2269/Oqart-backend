-- Drop triggers
DROP TRIGGER IF EXISTS trigger_ensure_one_default_address ON addresses;
DROP TRIGGER IF EXISTS trigger_update_cart_totals_on_item_change ON cart_items;
DROP TRIGGER IF EXISTS trigger_calculate_cart_item_subtotal ON cart_items;
DROP TRIGGER IF EXISTS update_wishlists_updated_at ON wishlists;
DROP TRIGGER IF EXISTS update_cart_items_updated_at ON cart_items;
DROP TRIGGER IF EXISTS update_carts_updated_at ON carts;
DROP TRIGGER IF EXISTS update_addresses_updated_at ON addresses;

-- Drop functions
DROP FUNCTION IF EXISTS ensure_one_default_address();
DROP FUNCTION IF EXISTS update_cart_totals();
DROP FUNCTION IF EXISTS calculate_cart_item_subtotal();

-- Drop tables
DROP TABLE IF EXISTS wishlist_items;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS addresses;

-- Drop ENUM types
DROP TYPE IF EXISTS address_type;
