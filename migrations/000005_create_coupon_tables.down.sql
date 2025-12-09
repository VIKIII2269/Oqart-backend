-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_coupon_status ON coupons;
DROP TRIGGER IF EXISTS update_coupons_updated_at ON coupons;

-- Drop functions
DROP FUNCTION IF EXISTS update_coupon_status();

-- Drop tables
DROP TABLE IF EXISTS coupon_user_restrictions;
DROP TABLE IF EXISTS coupon_usage;
DROP TABLE IF EXISTS coupons;

-- Drop ENUM types
DROP TYPE IF EXISTS coupon_user_type;
DROP TYPE IF EXISTS coupon_status;
DROP TYPE IF EXISTS coupon_type;
