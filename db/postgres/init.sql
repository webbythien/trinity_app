-- Enums
CREATE TYPE campaign_user_type AS ENUM ('internal', 'external', 'both');
CREATE TYPE package_type AS ENUM ('bronze', 'silver', 'gold');
CREATE TYPE voucher_status AS ENUM ('active', 'used', 'expired');
CREATE TYPE subscription_status AS ENUM ('active', 'expired', 'cancelled', 'pending');
CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'failed', 'refunded');
CREATE TYPE discount_type AS ENUM ('percentage', 'fixed');
CREATE TYPE discount_entity_type AS ENUM ('package', 'tour', 'merchandise');

-- Entity types management
CREATE TABLE entity_types (
    id SERIAL PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL UNIQUE,
    table_name VARCHAR(50) NOT NULL,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Entity views tracking
CREATE TABLE entity_views (
    id SERIAL PRIMARY KEY,
    entity_type_id INTEGER REFERENCES entity_types(id),
    view_name VARCHAR(100) NOT NULL,
    view_query TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Roles
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    role_id INTEGER REFERENCES roles(id),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Packages
CREATE TABLE packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    package_type package_type NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    description TEXT,
    duration_months INTEGER NOT NULL,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO packages (name, package_type, price, description, duration_months, status) VALUES
('Basic Package', 'bronze', 19.99, 'A basic subscription package for beginners.', 1, true),
('Standard Package', 'silver', 49.99, 'A standard package for regular users.', 6, true),
('Premium Package', 'gold', 99.99, 'A premium package for advanced users.', 12, true),
('Trial Package', 'bronze', 0.00, 'A free trial package for new users.', 1, true),
('Professional Package', 'silver', 69.99, 'A professional package for business users.', 12, false);


-- Marketing platforms
CREATE TABLE platforms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO platforms (name, status) VALUES
('Facebook', true),
('Google', true),
('Twitter', true),
('TikTok', true),
('YouTube', true),
('Instagram', true),
('LinkedIn', true),
('Snapchat', true),
('Pinterest', true),
('Reddit', true);


-- Promotional campaigns
CREATE TABLE promotional_campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    discount_type discount_type NOT NULL,
    discount_value DECIMAL(10,2) NOT NULL,
    max_discount_amount DECIMAL(10,2),
    user_type campaign_user_type NOT NULL,
    max_vouchers INTEGER NOT NULL,
    remaining_vouchers INTEGER NOT NULL,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Campaign entities junction
CREATE TABLE campaign_entities (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES promotional_campaigns(id),
    entity_type discount_entity_type NOT NULL,
    entity_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(campaign_id, entity_type, entity_id)
);

-- Campaign platform limits
CREATE TABLE campaign_platform_limits (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES promotional_campaigns(id),
    platform_id INTEGER REFERENCES platforms(id),
    voucher_limit INTEGER DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(campaign_id, platform_id)
);

-- Campaign tracking
CREATE TABLE campaign_tracking (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES promotional_campaigns(id),
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    referrer_url TEXT,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    converted BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vouchers
CREATE TABLE vouchers (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES promotional_campaigns(id),
    user_id INTEGER REFERENCES users(id),
    tracking_id INTEGER REFERENCES campaign_tracking(id),
    code VARCHAR(50) NOT NULL UNIQUE,
    discount_amount DECIMAL(10,2) NOT NULL,
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL,
    valid_until TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    status voucher_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Subscriptions
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    package_id INTEGER REFERENCES packages(id),
    voucher_id INTEGER REFERENCES vouchers(id),
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    original_price DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) NOT NULL,
    final_price DECIMAL(10,2) NOT NULL,
    status subscription_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default data
INSERT INTO entity_types (type_name, table_name) VALUES
('package', 'packages'),
('tour', 'tours'), -- easy expand user voucher for future
('merchandise', 'merchandise');

INSERT INTO roles (name) VALUES
('admin'),
('user');

-- Create views
CREATE VIEW active_packages AS
SELECT id, name, package_type, price
FROM packages
WHERE status = true;

INSERT INTO entity_views (entity_type_id, view_name, view_query)
SELECT id, 'active_packages', 
'SELECT id, name, package_type, price FROM packages WHERE status = true'
FROM entity_types 
WHERE type_name = 'package';