-- Create ENUM types for notifications
CREATE TYPE notification_type AS ENUM ('order', 'payment', 'return', 'review', 'vendor', 'system', 'marketing', 'general');
CREATE TYPE notification_channel AS ENUM ('email', 'sms', 'push', 'in_app');
CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'failed', 'read');

-- Notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type notification_type NOT NULL,
    channel notification_channel NOT NULL,
    status notification_status NOT NULL DEFAULT 'pending',

    -- Content
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,

    -- Links
    action_url VARCHAR(500),
    image_url VARCHAR(500),

    -- Status flags
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,

    -- Scheduling
    scheduled_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,

    -- Error tracking
    error_message TEXT,
    retry_count INT DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Audit logs table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Search logs table
CREATE TABLE search_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    session_id VARCHAR(255),
    query VARCHAR(500) NOT NULL,
    results_count INT DEFAULT 0,
    filters JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Serviceable pincodes table
CREATE TABLE serviceable_pincodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pincode VARCHAR(10) NOT NULL UNIQUE,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    district VARCHAR(100),
    region VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    is_cod_available BOOLEAN DEFAULT TRUE,
    standard_delivery_days INT DEFAULT 3,
    express_delivery_days INT DEFAULT 1,
    shipping_charge DECIMAL(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Newsletter subscribers table
CREATE TABLE newsletter_subscribers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    subscribed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    unsubscribed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Pages table (for CMS)
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    excerpt TEXT,
    status VARCHAR(50) DEFAULT 'published',
    is_active BOOLEAN DEFAULT TRUE,
    meta_title VARCHAR(255),
    meta_description TEXT,
    meta_keywords TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Banners table
CREATE TABLE banners (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(500) NOT NULL,
    mobile_image_url VARCHAR(500),
    link_url VARCHAR(500),
    position VARCHAR(50) NOT NULL DEFAULT 'home_slider',
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    click_count INT DEFAULT 0,
    view_count INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- FAQ table
CREATE TABLE faqs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category VARCHAR(100) NOT NULL,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Contact form submissions table
CREATE TABLE contact_submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    ip_address VARCHAR(45),
    status VARCHAR(50) DEFAULT 'new',
    responded_at TIMESTAMP WITH TIME ZONE,
    responded_by UUID REFERENCES users(id),
    response TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Testimonials table
CREATE TABLE testimonials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    designation VARCHAR(255),
    company VARCHAR(255),
    image_url VARCHAR(500),
    text TEXT NOT NULL,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    is_active BOOLEAN DEFAULT TRUE,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Settings table (for system-wide settings)
CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(255) NOT NULL UNIQUE,
    value TEXT,
    type VARCHAR(50) DEFAULT 'string',
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for notifications
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;

-- Indexes for audit_logs
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Indexes for search_logs
CREATE INDEX idx_search_logs_user_id ON search_logs(user_id);
CREATE INDEX idx_search_logs_query ON search_logs(query);
CREATE INDEX idx_search_logs_created_at ON search_logs(created_at);

-- Indexes for serviceable_pincodes
CREATE INDEX idx_serviceable_pincodes_pincode ON serviceable_pincodes(pincode);
CREATE INDEX idx_serviceable_pincodes_city ON serviceable_pincodes(city);
CREATE INDEX idx_serviceable_pincodes_state ON serviceable_pincodes(state);
CREATE INDEX idx_serviceable_pincodes_is_active ON serviceable_pincodes(is_active);

-- Indexes for newsletter_subscribers
CREATE INDEX idx_newsletter_subscribers_email ON newsletter_subscribers(email);
CREATE INDEX idx_newsletter_subscribers_status ON newsletter_subscribers(status);

-- Indexes for pages
CREATE INDEX idx_pages_slug ON pages(slug);
CREATE INDEX idx_pages_status ON pages(status);
CREATE INDEX idx_pages_is_active ON pages(is_active);

-- Indexes for banners
CREATE INDEX idx_banners_position ON banners(position);
CREATE INDEX idx_banners_is_active ON banners(is_active);
CREATE INDEX idx_banners_dates ON banners(start_date, end_date);

-- Indexes for faqs
CREATE INDEX idx_faqs_category ON faqs(category);
CREATE INDEX idx_faqs_is_active ON faqs(is_active);

-- Indexes for contact_submissions
CREATE INDEX idx_contact_submissions_email ON contact_submissions(email);
CREATE INDEX idx_contact_submissions_status ON contact_submissions(status);
CREATE INDEX idx_contact_submissions_created_at ON contact_submissions(created_at);

-- Indexes for testimonials
CREATE INDEX idx_testimonials_is_active ON testimonials(is_active);
CREATE INDEX idx_testimonials_user_id ON testimonials(user_id);

-- Indexes for settings
CREATE INDEX idx_settings_key ON settings(key);
CREATE INDEX idx_settings_is_public ON settings(is_public);

-- Add triggers for updated_at
CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_serviceable_pincodes_updated_at BEFORE UPDATE ON serviceable_pincodes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pages_updated_at BEFORE UPDATE ON pages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_banners_updated_at BEFORE UPDATE ON banners
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_faqs_updated_at BEFORE UPDATE ON faqs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_testimonials_updated_at BEFORE UPDATE ON testimonials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_settings_updated_at BEFORE UPDATE ON settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
