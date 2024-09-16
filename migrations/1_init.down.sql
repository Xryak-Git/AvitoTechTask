DROP TABLE IF EXISTS tender;

DROP TYPE IF EXISTS tender_status;

DROP TYPE IF EXISTS service_type;

DROP TABLE IF EXISTS organization_responsible;

DROP TABLE IF EXISTS employee;

DROP TABLE IF EXISTS bid;

DROP TYPE IF EXISTS bid_status;

DROP TYPE IF EXISTS author_type;

DROP TABLE IF EXISTS bid_review;

DROP TABLE IF EXISTS organization;

DROP TYPE IF EXISTS organization_type;

DROP TABLE IF EXISTS tender_versions;

DROP FUNCTION IF EXISTS save_and_increment_tender_version;

DROP TRIGGER IF EXISTS tender_update_trigger ON tender;

DROP TABLE IF EXISTS bid_bidReview;

DROP TABLE IF EXISTS bid_versions;

DROP FUNCTION IF EXISTS save_and_increment_bid_version;

DROP TRIGGER IF EXISTS bid_update_trigger ON bid;