-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_vendor_order_stats ON orders;
DROP TRIGGER IF EXISTS trigger_create_order_status_history ON orders;
DROP TRIGGER IF EXISTS trigger_generate_invoice_number ON invoices;
DROP TRIGGER IF EXISTS trigger_generate_order_number ON orders;
DROP TRIGGER IF EXISTS update_invoices_updated_at ON invoices;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;

-- Drop functions
DROP FUNCTION IF EXISTS update_vendor_order_stats();
DROP FUNCTION IF EXISTS create_order_status_history();
DROP FUNCTION IF EXISTS generate_invoice_number();
DROP FUNCTION IF EXISTS generate_order_number();

-- Drop sequences
DROP SEQUENCE IF EXISTS invoice_number_seq;
DROP SEQUENCE IF EXISTS order_number_seq;

-- Drop tables
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS payment_transactions;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- Drop ENUM types
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS order_status;
