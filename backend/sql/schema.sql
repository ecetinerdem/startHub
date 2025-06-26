-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table for role-based access
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('starthub', 'investor', 'donator', 'collaborator')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Starthubs table (now with image_url field)
CREATE TABLE IF NOT EXISTS starthubs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    location TEXT,
    team_size INT,
    url TEXT,
    email TEXT UNIQUE NOT NULL,
    image_url TEXT,
    join_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL  -- Add this line
);



-- Many-to-many: Starthub <-> Category
CREATE TABLE IF NOT EXISTS starthub_categories (
    starthub_id UUID REFERENCES starthubs(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (starthub_id, category_id)
);

-- Self-referencing many-to-many: Starthub collaborations
CREATE TABLE IF NOT EXISTS starthub_collaborations (
    starthub_id UUID REFERENCES starthubs(id) ON DELETE CASCADE,
    collaborator_id UUID REFERENCES starthubs(id) ON DELETE CASCADE,
    PRIMARY KEY (starthub_id, collaborator_id)
);

-- External collaborators not in starthubs
CREATE TABLE IF NOT EXISTS external_collaborators (
    id SERIAL PRIMARY KEY,
    starthub_id UUID REFERENCES starthubs(id) ON DELETE CASCADE,
    name TEXT NOT NULL
);


