-- Initial schema for credfolio2
-- Creates core tables: users, files, reference_letters

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table: stores user accounts
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for email lookups during authentication
CREATE INDEX idx_users_email ON users(email);

-- Files table: tracks uploaded files stored in MinIO
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    storage_key VARCHAR(512) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for looking up files by user
CREATE INDEX idx_files_user_id ON files(user_id);

-- Reference letters table: stores extracted data from reference letters
CREATE TABLE reference_letters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_id UUID REFERENCES files(id) ON DELETE SET NULL,
    title VARCHAR(255),
    author_name VARCHAR(255),
    author_title VARCHAR(255),
    organization VARCHAR(255),
    date_written DATE,
    raw_text TEXT,
    extracted_data JSONB,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for looking up reference letters by user
CREATE INDEX idx_reference_letters_user_id ON reference_letters(user_id);

-- Index for status filtering
CREATE INDEX idx_reference_letters_status ON reference_letters(status);

-- Trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply update trigger to users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply update trigger to reference_letters table
CREATE TRIGGER update_reference_letters_updated_at
    BEFORE UPDATE ON reference_letters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
