CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Определение ENUM типов
CREATE TYPE service_type AS ENUM (
    'CONSTRUCTION',
    'DELIVERY',
    'MANUFACTURING'
    );

CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
    );

CREATE TYPE bid_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CANCELED',
    'APPROVED',
    'REJECTED'
    );

CREATE TYPE author_type AS ENUM (
    'ORGANIZATION',
    'USER'
    );

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );


CREATE TABLE IF NOT EXISTS organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tender (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    service_type service_type,
    status tender_status DEFAULT 'CREATED',
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    status bid_status DEFAULT 'CREATED',
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_type author_type,
    author_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid_review (
    id UUID PRIMARY KEY,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE tender_versions (
    id SERIAL PRIMARY KEY,
    tender_id uuid REFERENCES tender(id),
    name TEXT,
    description TEXT,
    service_type service_type,
    status TEXT,
    organization_id uuid REFERENCES organization(id),
    version INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION save_and_increment_tender_version() RETURNS TRIGGER AS $$
BEGIN
    -- Вставляем старую версию в таблицу tender_versions
    INSERT INTO tender_versions (tender_id, name, description, service_type, status, organization_id, version)
    VALUES (OLD.id, OLD.name, OLD.description, OLD.service_type, OLD.status, OLD.organization_id, OLD.version);

    NEW.version := OLD.version + 1;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER tender_update_trigger
    BEFORE UPDATE ON tender
    FOR EACH ROW
EXECUTE FUNCTION save_and_increment_tender_version();

CREATE TABLE bid_bidReview (
    bid_id UUID NOT NULL,
    bid_review_id UUID NOT NULL,
    PRIMARY KEY (bid_id, bid_review_id),
    FOREIGN KEY (bid_id) REFERENCES bid(id),
    FOREIGN KEY (bid_review_id) REFERENCES bid_review(id)
);

CREATE TABLE bid_versions (
    id SERIAL PRIMARY KEY,
    bid_id uuid REFERENCES bid(id),
    name TEXT,
    description TEXT,
    status bid_status,
    tender_id uuid REFERENCES tender(id),
    author_type author_type,
    author_id uuid REFERENCES employee(id),
    version INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION save_and_increment_bid_version() RETURNS TRIGGER AS $$
BEGIN

    INSERT INTO bid_versions (bid_id, name, description, status, tender_id, author_type, author_id, version)
    VALUES (OLD.id, OLD.name, OLD.description, OLD.status, OLD.tender_id, OLD.author_type, OLD.author_id,OLD.version);

    NEW.version := OLD.version + 1;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER bid_update_trigger
    BEFORE UPDATE ON bid
    FOR EACH ROW
EXECUTE FUNCTION save_and_increment_bid_version();


CREATE TABLE decision (
    id SERIAL PRIMARY KEY,
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    decision_type bid_status
);

CREATE OR REPLACE FUNCTION update_bid_status()
    RETURNS TRIGGER AS $$
BEGIN

    DECLARE
        v_has_reject BOOLEAN;
        v_approve_count INT;
        v_quorum INT;
    BEGIN

        SELECT COUNT(*) > 0 INTO v_has_reject
        FROM decision
        WHERE bid_id = NEW.bid_id
          AND decision_type = 'REJECTED'::bid_status;


        SELECT COUNT(*) INTO v_approve_count
        FROM decision
        WHERE bid_id = NEW.bid_id
          AND decision_type = 'APPROVED'::bid_status;


        SELECT LEAST(3, COUNT(*)) INTO v_quorum
        FROM organization_responsible or_res
                 JOIN bid b ON b.author_id = or_res.user_id
        WHERE b.id = NEW.bid_id;


        IF v_has_reject THEN
            UPDATE bid
            SET status = 'REJECTED'::bid_status
            WHERE id = NEW.bid_id;
        ELSIF v_approve_count >= v_quorum THEN
            UPDATE bid
            SET status = 'APPROVED'::bid_status
            WHERE id = NEW.bid_id;
        END IF;
        RETURN NEW;
    END;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER decision_status_trigger
    AFTER INSERT OR UPDATE ON decision
    FOR EACH ROW
EXECUTE FUNCTION update_bid_status();
