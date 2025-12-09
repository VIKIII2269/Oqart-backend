-- Drop triggers
DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;
DROP TRIGGER IF EXISTS update_testimonials_updated_at ON testimonials;
DROP TRIGGER IF EXISTS update_faqs_updated_at ON faqs;
DROP TRIGGER IF EXISTS update_banners_updated_at ON banners;
DROP TRIGGER IF EXISTS update_pages_updated_at ON pages;
DROP TRIGGER IF EXISTS update_serviceable_pincodes_updated_at ON serviceable_pincodes;
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;

-- Drop tables
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS testimonials;
DROP TABLE IF EXISTS contact_submissions;
DROP TABLE IF EXISTS faqs;
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS newsletter_subscribers;
DROP TABLE IF EXISTS serviceable_pincodes;
DROP TABLE IF EXISTS search_logs;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS notifications;

-- Drop ENUM types
DROP TYPE IF EXISTS notification_status;
DROP TYPE IF EXISTS notification_channel;
DROP TYPE IF EXISTS notification_type;
