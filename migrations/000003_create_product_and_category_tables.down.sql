-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_product_discount ON products;
DROP TRIGGER IF EXISTS trigger_update_product_stock_status ON products;
DROP TRIGGER IF EXISTS update_product_variants_updated_at ON product_variants;
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TRIGGER IF EXISTS update_certifications_updated_at ON certifications;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;

-- Drop functions
DROP FUNCTION IF EXISTS update_product_discount();
DROP FUNCTION IF EXISTS update_product_stock_status();

-- Drop tables
DROP TABLE IF EXISTS product_specifications;
DROP TABLE IF EXISTS product_tag_mapping;
DROP TABLE IF EXISTS product_certifications;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS product_tags;
DROP TABLE IF EXISTS certifications;
DROP TABLE IF EXISTS categories;

-- Drop ENUM types
DROP TYPE IF EXISTS stock_status;
DROP TYPE IF EXISTS product_status;
