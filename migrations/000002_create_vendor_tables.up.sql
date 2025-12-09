-- Create ENUM types for vendors
CREATE TYPE vendor_status AS ENUM ('pending', 'under_review', 'approved', 'rejected', 'suspended', 'inactive');
CREATE TYPE vendor_business_type AS ENUM ('individual', 'proprietorship', 'partnership', 'private_limited', 'public_limited', 'llp');
CREATE TYPE document_type AS ENUM ('gst_certificate', 'pan_card', 'address_proof', 'bank_statement', 'cancelled_cheque', 'fssai_license', 'organic_certificate');
CREATE TYPE document_status AS ENUM ('pending', 'under_review', 'approved', 'rejected');

-- Vendors table
CREATE TABLE vendors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    business_name VARCHAR(255) NOT NULL,
    business_type vendor_business_type NOT NULL,
    gstin VARCHAR(15) UNIQUE,
    pan VARCHAR(10),
    fssai_number VARCHAR(14),
    status vendor_status NOT NULL DEFAULT 'pending',
    onboarding_step INT DEFAULT 1,
    rejection_reason TEXT,
    approved_at TIMESTAMP WITH TIME ZONE,
    approved_by UUID REFERENCES users(id),
    rating DECIMAL(3, 2) DEFAULT 0.00,
    total_reviews INT DEFAULT 0,
    total_orders INT DEFAULT 0,
    total_sales DECIMAL(15, 2) DEFAULT 0.00,
    commission_rate DECIMAL(5, 2) DEFAULT 10.00,
    is_verified BOOLEAN DEFAULT FALSE,
    is_featured BOOLEAN DEFAULT FALSE,
    description TEXT,
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vendor addresses table
CREATE TABLE vendor_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL DEFAULT 'business',
    street_address VARCHAR(255) NOT NULL,
    landmark VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    pincode VARCHAR(10) NOT NULL,
    country VARCHAR(100) DEFAULT 'India',
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vendor documents table
CREATE TABLE vendor_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    type document_type NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    status document_status NOT NULL DEFAULT 'pending',
    extracted_data JSONB,
    rejection_reason TEXT,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vendor bank details table
CREATE TABLE vendor_bank_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL UNIQUE REFERENCES vendors(id) ON DELETE CASCADE,
    account_holder_name VARCHAR(255) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    ifsc_code VARCHAR(11) NOT NULL,
    bank_name VARCHAR(255) NOT NULL,
    branch_name VARCHAR(255),
    account_type VARCHAR(50) DEFAULT 'savings',
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vendor contact persons table
CREATE TABLE vendor_contact_persons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    designation VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vendor business hours table
CREATE TABLE vendor_business_hours (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    is_open BOOLEAN DEFAULT TRUE,
    opening_time TIME,
    closing_time TIME,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(vendor_id, day_of_week)
);

-- Indexes for vendors
CREATE INDEX idx_vendors_user_id ON vendors(user_id);
CREATE INDEX idx_vendors_status ON vendors(status);
CREATE INDEX idx_vendors_gstin ON vendors(gstin);
CREATE INDEX idx_vendors_rating ON vendors(rating);
CREATE INDEX idx_vendors_is_featured ON vendors(is_featured);
CREATE INDEX idx_vendors_created_at ON vendors(created_at);

-- Indexes for vendor_addresses
CREATE INDEX idx_vendor_addresses_vendor_id ON vendor_addresses(vendor_id);
CREATE INDEX idx_vendor_addresses_pincode ON vendor_addresses(pincode);
CREATE INDEX idx_vendor_addresses_city ON vendor_addresses(city);
CREATE INDEX idx_vendor_addresses_state ON vendor_addresses(state);

-- Indexes for vendor_documents
CREATE INDEX idx_vendor_documents_vendor_id ON vendor_documents(vendor_id);
CREATE INDEX idx_vendor_documents_type ON vendor_documents(type);
CREATE INDEX idx_vendor_documents_status ON vendor_documents(status);

-- Indexes for vendor_bank_details
CREATE INDEX idx_vendor_bank_details_vendor_id ON vendor_bank_details(vendor_id);

-- Indexes for vendor_contact_persons
CREATE INDEX idx_vendor_contact_persons_vendor_id ON vendor_contact_persons(vendor_id);

-- Add triggers for updated_at
CREATE TRIGGER update_vendors_updated_at BEFORE UPDATE ON vendors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vendor_addresses_updated_at BEFORE UPDATE ON vendor_addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vendor_documents_updated_at BEFORE UPDATE ON vendor_documents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vendor_bank_details_updated_at BEFORE UPDATE ON vendor_bank_details
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vendor_contact_persons_updated_at BEFORE UPDATE ON vendor_contact_persons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vendor_business_hours_updated_at BEFORE UPDATE ON vendor_business_hours
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
