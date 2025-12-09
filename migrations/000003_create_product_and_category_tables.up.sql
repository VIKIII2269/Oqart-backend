-- Create ENUM types for products
CREATE TYPE product_status AS ENUM ('draft', 'pending', 'approved', 'rejected', 'active', 'inactive', 'out_of_stock');
CREATE TYPE stock_status AS ENUM ('in_stock', 'low_stock', 'out_of_stock', 'discontinued');

-- Categories table (hierarchical)
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    image_url VARCHAR(500),
    icon_url VARCHAR(500),
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    meta_title VARCHAR(255),
    meta_description TEXT,
    meta_keywords TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Certifications table
CREATE TABLE certifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    icon_url VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Product tags table
CREATE TABLE product_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Products table
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    sku VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    short_description TEXT,
    status product_status NOT NULL DEFAULT 'draft',
    stock_status stock_status NOT NULL DEFAULT 'in_stock',

    -- Pricing
    price DECIMAL(15, 2) NOT NULL,
    mrp DECIMAL(15, 2) NOT NULL,
    discount_percentage DECIMAL(5, 2) DEFAULT 0.00,
    cost_price DECIMAL(15, 2),

    -- Stock management
    stock_quantity INT DEFAULT 0,
    low_stock_threshold INT DEFAULT 10,
    allow_backorder BOOLEAN DEFAULT FALSE,

    -- Product details
    weight DECIMAL(10, 2),
    weight_unit VARCHAR(20) DEFAULT 'kg',
    dimensions_length DECIMAL(10, 2),
    dimensions_width DECIMAL(10, 2),
    dimensions_height DECIMAL(10, 2),
    dimensions_unit VARCHAR(20) DEFAULT 'cm',

    -- Additional info
    brand VARCHAR(255),
    manufacturer VARCHAR(255),
    country_of_origin VARCHAR(100) DEFAULT 'India',
    shelf_life VARCHAR(100),
    ingredients TEXT,
    nutritional_info JSONB,
    usage_instructions TEXT,
    storage_instructions TEXT,

    -- Features
    is_featured BOOLEAN DEFAULT FALSE,
    is_bestseller BOOLEAN DEFAULT FALSE,
    is_new_arrival BOOLEAN DEFAULT FALSE,
    is_organic BOOLEAN DEFAULT TRUE,
    is_vegan BOOLEAN DEFAULT FALSE,
    is_gluten_free BOOLEAN DEFAULT FALSE,

    -- SEO
    meta_title VARCHAR(255),
    meta_description TEXT,
    meta_keywords TEXT,

    -- Stats
    views_count INT DEFAULT 0,
    sales_count INT DEFAULT 0,
    rating DECIMAL(3, 2) DEFAULT 0.00,
    review_count INT DEFAULT 0,

    -- Timestamps
    approved_at TIMESTAMP WITH TIME ZONE,
    approved_by UUID REFERENCES users(id),
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Product images table
CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    display_order INT DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Product variants table
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(15, 2) NOT NULL,
    mrp DECIMAL(15, 2) NOT NULL,
    stock_quantity INT DEFAULT 0,
    weight DECIMAL(10, 2),
    weight_unit VARCHAR(20) DEFAULT 'kg',
    attributes JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Product certifications mapping table
CREATE TABLE product_certifications (
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    certification_id UUID NOT NULL REFERENCES certifications(id) ON DELETE CASCADE,
    certificate_number VARCHAR(255),
    valid_from DATE,
    valid_until DATE,
    certificate_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, certification_id)
);

-- Product tags mapping table
CREATE TABLE product_tag_mapping (
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES product_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, tag_id)
);

-- Product specifications table
CREATE TABLE product_specifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for categories
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_is_active ON categories(is_active);
CREATE INDEX idx_categories_is_featured ON categories(is_featured);

-- Indexes for products
CREATE INDEX idx_products_vendor_id ON products(vendor_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_stock_status ON products(stock_status);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_is_featured ON products(is_featured);
CREATE INDEX idx_products_rating ON products(rating);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_products_vendor_category ON products(vendor_id, category_id);
CREATE INDEX idx_products_status_featured ON products(status, is_featured);

-- Full-text search index
CREATE INDEX idx_products_search ON products USING gin(to_tsvector('english', name || ' ' || COALESCE(description, '')));

-- Indexes for product_images
CREATE INDEX idx_product_images_product_id ON product_images(product_id);
CREATE INDEX idx_product_images_is_primary ON product_images(is_primary);

-- Indexes for product_variants
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_sku ON product_variants(sku);

-- Indexes for product_tags
CREATE INDEX idx_product_tags_slug ON product_tags(slug);

-- Indexes for product_specifications
CREATE INDEX idx_product_specifications_product_id ON product_specifications(product_id);

-- Add triggers for updated_at
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_certifications_updated_at BEFORE UPDATE ON certifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_product_variants_updated_at BEFORE UPDATE ON product_variants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to automatically update stock status
CREATE OR REPLACE FUNCTION update_product_stock_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.stock_quantity = 0 THEN
        NEW.stock_status = 'out_of_stock';
    ELSIF NEW.stock_quantity <= NEW.low_stock_threshold THEN
        NEW.stock_status = 'low_stock';
    ELSE
        NEW.stock_status = 'in_stock';
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_product_stock_status
    BEFORE INSERT OR UPDATE OF stock_quantity ON products
    FOR EACH ROW EXECUTE FUNCTION update_product_stock_status();

-- Function to update discount percentage
CREATE OR REPLACE FUNCTION update_product_discount()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.mrp > 0 THEN
        NEW.discount_percentage = ROUND(((NEW.mrp - NEW.price) / NEW.mrp * 100)::numeric, 2);
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_product_discount
    BEFORE INSERT OR UPDATE OF price, mrp ON products
    FOR EACH ROW EXECUTE FUNCTION update_product_discount();
