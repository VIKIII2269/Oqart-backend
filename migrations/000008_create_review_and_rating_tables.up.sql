-- Create ENUM types for reviews
CREATE TYPE review_status AS ENUM ('pending', 'approved', 'rejected', 'flagged');

-- Reviews table
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    order_item_id UUID REFERENCES order_items(id) ON DELETE SET NULL,

    -- Rating and review
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    text TEXT NOT NULL,
    images JSONB,

    -- Status
    status review_status NOT NULL DEFAULT 'approved',

    -- Helpful votes
    helpful_count INT DEFAULT 0,
    not_helpful_count INT DEFAULT 0,

    -- Moderation
    is_verified_purchase BOOLEAN DEFAULT FALSE,
    flagged_count INT DEFAULT 0,
    moderated_by UUID REFERENCES users(id),
    moderated_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,

    -- Vendor response
    vendor_response TEXT,
    vendor_response_by UUID REFERENCES users(id),
    vendor_response_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Review votes table (helpful/not helpful)
CREATE TABLE review_votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    review_id UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_helpful BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(review_id, user_id)
);

-- Review flags table
CREATE TABLE review_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    review_id UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(255) NOT NULL,
    details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(review_id, user_id)
);

-- Vendor ratings table (overall vendor ratings)
CREATE TABLE vendor_ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,

    -- Ratings
    product_quality INT NOT NULL CHECK (product_quality >= 1 AND product_quality <= 5),
    packaging INT NOT NULL CHECK (packaging >= 1 AND packaging <= 5),
    delivery_speed INT NOT NULL CHECK (delivery_speed >= 1 AND delivery_speed <= 5),
    customer_service INT NOT NULL CHECK (customer_service >= 1 AND customer_service <= 5),
    overall_rating DECIMAL(3, 2) NOT NULL,

    -- Review
    review_text TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(vendor_id, user_id, order_id)
);

-- Indexes for reviews
CREATE INDEX idx_reviews_product_id ON reviews(product_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_order_id ON reviews(order_id);
CREATE INDEX idx_reviews_status ON reviews(status);
CREATE INDEX idx_reviews_rating ON reviews(rating);
CREATE INDEX idx_reviews_created_at ON reviews(created_at);
CREATE INDEX idx_reviews_product_status ON reviews(product_id, status);

-- Indexes for review_votes
CREATE INDEX idx_review_votes_review_id ON review_votes(review_id);
CREATE INDEX idx_review_votes_user_id ON review_votes(user_id);

-- Indexes for review_flags
CREATE INDEX idx_review_flags_review_id ON review_flags(review_id);
CREATE INDEX idx_review_flags_user_id ON review_flags(user_id);

-- Indexes for vendor_ratings
CREATE INDEX idx_vendor_ratings_vendor_id ON vendor_ratings(vendor_id);
CREATE INDEX idx_vendor_ratings_user_id ON vendor_ratings(user_id);
CREATE INDEX idx_vendor_ratings_order_id ON vendor_ratings(order_id);

-- Add triggers for updated_at
CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vendor_ratings_updated_at BEFORE UPDATE ON vendor_ratings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to update product rating
CREATE OR REPLACE FUNCTION update_product_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE products
    SET rating = (
        SELECT ROUND(AVG(rating)::numeric, 2)
        FROM reviews
        WHERE product_id = NEW.product_id AND status = 'approved'
    ),
    review_count = (
        SELECT COUNT(*)
        FROM reviews
        WHERE product_id = NEW.product_id AND status = 'approved'
    )
    WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_product_rating
    AFTER INSERT OR UPDATE OF rating, status ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_product_rating();

-- Function to update vendor rating
CREATE OR REPLACE FUNCTION update_vendor_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE vendors
    SET rating = (
        SELECT ROUND(AVG(overall_rating)::numeric, 2)
        FROM vendor_ratings
        WHERE vendor_id = NEW.vendor_id
    ),
    total_reviews = (
        SELECT COUNT(*)
        FROM vendor_ratings
        WHERE vendor_id = NEW.vendor_id
    )
    WHERE id = NEW.vendor_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_vendor_rating
    AFTER INSERT OR UPDATE OF overall_rating ON vendor_ratings
    FOR EACH ROW EXECUTE FUNCTION update_vendor_rating();

-- Function to update review helpful counts
CREATE OR REPLACE FUNCTION update_review_vote_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.is_helpful THEN
            UPDATE reviews SET helpful_count = helpful_count + 1 WHERE id = NEW.review_id;
        ELSE
            UPDATE reviews SET not_helpful_count = not_helpful_count + 1 WHERE id = NEW.review_id;
        END IF;
    ELSIF TG_OP = 'UPDATE' AND OLD.is_helpful != NEW.is_helpful THEN
        IF NEW.is_helpful THEN
            UPDATE reviews
            SET helpful_count = helpful_count + 1, not_helpful_count = not_helpful_count - 1
            WHERE id = NEW.review_id;
        ELSE
            UPDATE reviews
            SET helpful_count = helpful_count - 1, not_helpful_count = not_helpful_count + 1
            WHERE id = NEW.review_id;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.is_helpful THEN
            UPDATE reviews SET helpful_count = helpful_count - 1 WHERE id = OLD.review_id;
        ELSE
            UPDATE reviews SET not_helpful_count = not_helpful_count - 1 WHERE id = OLD.review_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_review_vote_counts
    AFTER INSERT OR UPDATE OR DELETE ON review_votes
    FOR EACH ROW EXECUTE FUNCTION update_review_vote_counts();

-- Function to update review flag count
CREATE OR REPLACE FUNCTION update_review_flag_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE reviews SET flagged_count = flagged_count + 1 WHERE id = NEW.review_id;
        -- Auto-flag for moderation if flagged more than 3 times
        UPDATE reviews SET status = 'flagged'
        WHERE id = NEW.review_id AND flagged_count >= 3 AND status = 'approved';
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE reviews SET flagged_count = flagged_count - 1 WHERE id = OLD.review_id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_review_flag_count
    AFTER INSERT OR DELETE ON review_flags
    FOR EACH ROW EXECUTE FUNCTION update_review_flag_count();
