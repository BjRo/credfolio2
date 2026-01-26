-- Seed demo user for testing
INSERT INTO users (id, email, password_hash, name)
VALUES ('00000000-0000-0000-0000-000000000001', 'demo@example.com', 'demo_hash', 'Demo User')
ON CONFLICT (id) DO NOTHING;
