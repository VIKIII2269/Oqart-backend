-- Drop triggers
DROP TRIGGER IF EXISTS update_vendor_business_hours_updated_at ON vendor_business_hours;
DROP TRIGGER IF EXISTS update_vendor_contact_persons_updated_at ON vendor_contact_persons;
DROP TRIGGER IF EXISTS update_vendor_bank_details_updated_at ON vendor_bank_details;
DROP TRIGGER IF EXISTS update_vendor_documents_updated_at ON vendor_documents;
DROP TRIGGER IF EXISTS update_vendor_addresses_updated_at ON vendor_addresses;
DROP TRIGGER IF EXISTS update_vendors_updated_at ON vendors;

-- Drop tables
DROP TABLE IF EXISTS vendor_business_hours;
DROP TABLE IF EXISTS vendor_contact_persons;
DROP TABLE IF EXISTS vendor_bank_details;
DROP TABLE IF EXISTS vendor_documents;
DROP TABLE IF EXISTS vendor_addresses;
DROP TABLE IF EXISTS vendors;

-- Drop ENUM types
DROP TYPE IF EXISTS document_status;
DROP TYPE IF EXISTS document_type;
DROP TYPE IF EXISTS vendor_business_type;
DROP TYPE IF EXISTS vendor_status;
