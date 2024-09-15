-- Создаем все необходимые сущности и типы
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE type serviceType_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);
CREATE type status_type AS ENUM (
    'Created',
    'Published',
    'Closed'
);
CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    serviceType serviceType_type,
    status status_type,
    organizationId UUID REFERENCES organization(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1 CHECK (version >= 1),
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE type authorType_type AS ENUM (
    'Organization',
    'User'
);

CREATE TABLE bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    status status_type,
    tenderId UUID REFERENCES tender(id) ON DELETE CASCADE,
    authorType authorType_type,
    authorId UUID REFERENCES employee(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1 CHECK (version >= 1),
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Вставляем данные 
INSERT INTO employee (id, username, first_name, last_name) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'test_user', 'test', 'user');

INSERT INTO organization (id, name, description, type) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Test org 1', 'string', 'IE');

INSERT INTO organization_responsible (organization_id, user_id) VALUES
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000');

INSERT INTO tender (id, name, description, serviceType, status, organizationId) VALUES
('c99d2f29-28d1-481a-9f38-3e9474991fe1', 'Tender test', 'Description for Tender test', 'Construction', 'Created', '550e8400-e29b-41d4-a716-446655440000');

INSERT INTO bid (name, description, status, tenderId, authorType, authorId, version)
VALUES ('Sample Bid', 'This is a sample bid description.', 'Created', 'c99d2f29-28d1-481a-9f38-3e9474991fe1', 'Organization', '550e8400-e29b-41d4-a716-446655440000', 1);

-- Извлекаем данные
-- SELECT * FROM employee;
-- SELECT * FROM organization;
-- SELECT * FROM organization_responsible;
-- SELECT * FROM tender;
-- SELECT * FROM bid;