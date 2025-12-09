-- Drop triggers
DROP TRIGGER IF EXISTS trigger_create_return_status_history ON returns;
DROP TRIGGER IF EXISTS trigger_generate_refund_number ON refunds;
DROP TRIGGER IF EXISTS trigger_generate_return_number ON returns;
DROP TRIGGER IF EXISTS update_refunds_updated_at ON refunds;
DROP TRIGGER IF EXISTS update_returns_updated_at ON returns;

-- Drop functions
DROP FUNCTION IF EXISTS create_return_status_history();
DROP FUNCTION IF EXISTS generate_refund_number();
DROP FUNCTION IF EXISTS generate_return_number();

-- Drop sequences
DROP SEQUENCE IF EXISTS refund_number_seq;
DROP SEQUENCE IF EXISTS return_number_seq;

-- Drop tables
DROP TABLE IF EXISTS return_status_history;
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS return_items;
DROP TABLE IF EXISTS returns;

-- Drop ENUM types
DROP TYPE IF EXISTS refund_method;
DROP TYPE IF EXISTS refund_status;
DROP TYPE IF EXISTS return_reason;
DROP TYPE IF EXISTS return_status;
