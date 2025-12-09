-- Create ENUM types for returns
CREATE TYPE return_status AS ENUM ('requested', 'pending', 'approved', 'rejected', 'picked_up', 'received', 'qc_passed', 'qc_failed', 'refunded', 'cancelled');
CREATE TYPE return_reason AS ENUM ('damaged', 'wrong_item', 'not_as_described', 'expired', 'missing_parts', 'quality_issue', 'changed_mind', 'other');
CREATE TYPE refund_status AS ENUM ('pending', 'processing', 'completed', 'failed', 'cancelled');
CREATE TYPE refund_method AS ENUM ('original', 'wallet', 'bank_transfer');

-- Returns table
CREATE TABLE returns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    return_number VARCHAR(50) NOT NULL UNIQUE,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE RESTRICT,
    status return_status NOT NULL DEFAULT 'requested',
    reason return_reason NOT NULL,
    reason_details TEXT,

    -- Return details
    return_amount DECIMAL(15, 2) NOT NULL,
    refund_amount DECIMAL(15, 2),
    shipping_refund DECIMAL(15, 2) DEFAULT 0.00,

    -- Images and proof
    images JSONB,

    -- Pickup details
    pickup_address JSONB,
    pickup_scheduled_at TIMESTAMP WITH TIME ZONE,
    pickup_completed_at TIMESTAMP WITH TIME ZONE,

    -- QC details
    qc_notes TEXT,
    qc_images JSONB,
    qc_passed BOOLEAN,
    qc_completed_at TIMESTAMP WITH TIME ZONE,

    -- Rejection
    rejection_reason TEXT,
    rejected_by UUID REFERENCES users(id),
    rejected_at TIMESTAMP WITH TIME ZONE,

    -- Approval
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Return items table
CREATE TABLE return_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    return_id UUID NOT NULL REFERENCES returns(id) ON DELETE CASCADE,
    order_item_id UUID NOT NULL REFERENCES order_items(id) ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,

    -- Item details snapshot
    product_name VARCHAR(255) NOT NULL,
    product_sku VARCHAR(100) NOT NULL,
    variant_name VARCHAR(255),

    -- Return details
    quantity INT NOT NULL CHECK (quantity > 0),
    price DECIMAL(15, 2) NOT NULL,
    subtotal DECIMAL(15, 2) NOT NULL,

    -- QC details
    qc_passed BOOLEAN,
    qc_notes TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Refunds table
CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    refund_number VARCHAR(50) NOT NULL UNIQUE,
    return_id UUID REFERENCES returns(id) ON DELETE SET NULL,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    payment_id UUID REFERENCES payments(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    status refund_status NOT NULL DEFAULT 'pending',
    method refund_method NOT NULL,

    -- Amounts
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'INR',

    -- Gateway details
    gateway VARCHAR(50),
    gateway_refund_id VARCHAR(255),
    gateway_response JSONB,

    -- Bank details (if refund_method = bank_transfer)
    bank_account_number VARCHAR(50),
    bank_ifsc VARCHAR(11),
    bank_account_holder VARCHAR(255),

    -- Processing details
    processed_by UUID REFERENCES users(id),
    processed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,

    -- Notes
    notes TEXT,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Return status history table
CREATE TABLE return_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    return_id UUID NOT NULL REFERENCES returns(id) ON DELETE CASCADE,
    status return_status NOT NULL,
    notes TEXT,
    changed_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for returns
CREATE INDEX idx_returns_return_number ON returns(return_number);
CREATE INDEX idx_returns_order_id ON returns(order_id);
CREATE INDEX idx_returns_user_id ON returns(user_id);
CREATE INDEX idx_returns_vendor_id ON returns(vendor_id);
CREATE INDEX idx_returns_status ON returns(status);
CREATE INDEX idx_returns_created_at ON returns(created_at);

-- Indexes for return_items
CREATE INDEX idx_return_items_return_id ON return_items(return_id);
CREATE INDEX idx_return_items_order_item_id ON return_items(order_item_id);
CREATE INDEX idx_return_items_product_id ON return_items(product_id);

-- Indexes for refunds
CREATE INDEX idx_refunds_refund_number ON refunds(refund_number);
CREATE INDEX idx_refunds_return_id ON refunds(return_id);
CREATE INDEX idx_refunds_order_id ON refunds(order_id);
CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_user_id ON refunds(user_id);
CREATE INDEX idx_refunds_status ON refunds(status);

-- Indexes for return_status_history
CREATE INDEX idx_return_status_history_return_id ON return_status_history(return_id);

-- Add triggers for updated_at
CREATE TRIGGER update_returns_updated_at BEFORE UPDATE ON returns
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_refunds_updated_at BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to generate return number
CREATE OR REPLACE FUNCTION generate_return_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.return_number IS NULL THEN
        NEW.return_number = 'RET' || TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDD') || LPAD(nextval('return_number_seq')::TEXT, 6, '0');
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create sequence for return numbers
CREATE SEQUENCE return_number_seq START 1;

CREATE TRIGGER trigger_generate_return_number
    BEFORE INSERT ON returns
    FOR EACH ROW EXECUTE FUNCTION generate_return_number();

-- Function to generate refund number
CREATE OR REPLACE FUNCTION generate_refund_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.refund_number IS NULL THEN
        NEW.refund_number = 'REF' || TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDD') || LPAD(nextval('refund_number_seq')::TEXT, 6, '0');
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create sequence for refund numbers
CREATE SEQUENCE refund_number_seq START 1;

CREATE TRIGGER trigger_generate_refund_number
    BEFORE INSERT ON refunds
    FOR EACH ROW EXECUTE FUNCTION generate_refund_number();

-- Function to automatically create return status history
CREATE OR REPLACE FUNCTION create_return_status_history()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' AND OLD.status != NEW.status) THEN
        INSERT INTO return_status_history (return_id, status, notes)
        VALUES (NEW.id, NEW.status, NULL);
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_create_return_status_history
    AFTER INSERT OR UPDATE OF status ON returns
    FOR EACH ROW EXECUTE FUNCTION create_return_status_history();
