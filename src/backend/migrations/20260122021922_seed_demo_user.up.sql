-- Seed demo user for testing
-- password_hash is bcrypt hash of 'demo123' for testing only
INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'demo@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mrq4H6E3JQKx6VLOVjHPb7dpqW8WJXW',
    'Demo User',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;
