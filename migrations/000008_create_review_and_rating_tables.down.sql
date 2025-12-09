-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_review_flag_count ON review_flags;
DROP TRIGGER IF EXISTS trigger_update_review_vote_counts ON review_votes;
DROP TRIGGER IF EXISTS trigger_update_vendor_rating ON vendor_ratings;
DROP TRIGGER IF EXISTS trigger_update_product_rating ON reviews;
DROP TRIGGER IF EXISTS update_vendor_ratings_updated_at ON vendor_ratings;
DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;

-- Drop functions
DROP FUNCTION IF EXISTS update_review_flag_count();
DROP FUNCTION IF EXISTS update_review_vote_counts();
DROP FUNCTION IF EXISTS update_vendor_rating();
DROP FUNCTION IF EXISTS update_product_rating();

-- Drop tables
DROP TABLE IF EXISTS vendor_ratings;
DROP TABLE IF EXISTS review_flags;
DROP TABLE IF EXISTS review_votes;
DROP TABLE IF EXISTS reviews;

-- Drop ENUM types
DROP TYPE IF EXISTS review_status;
