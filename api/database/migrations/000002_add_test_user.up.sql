-- Add test data for development and testing

-- Insert a test university
INSERT INTO universities (id, name_en, name_fa, is_active)
VALUES (
    '00000000-0000-4000-a000-000000000001',
    'Test University',
    'دانشگاه تست',
    true
) ON CONFLICT DO NOTHING;

-- Insert a test faculty
INSERT INTO faculties (id, university_id, name_en, name_fa, short_code, is_active)
VALUES (
    '00000000-0000-4000-a000-000000000002',
    '00000000-0000-4000-a000-000000000001',
    'Test Faculty',
    'دانشکده تست',
    'TEST',
    true
) ON CONFLICT DO NOTHING;

-- Insert a test user with email_verified set to true
INSERT INTO users (
    id,
    email,
    password_hash,
    student_id,
    first_name,
    last_name,
    university_id,
    faculty_id,
    gender,
    email_verified,
    is_admin
)
VALUES (
    '00000000-0000-4000-a000-000000000003',
    'test@example.com',
    -- This is the hash for 'strongpassword'
    '$2a$10$WUwgBnmz2WXqGTX8k/SrGeQS9J82tkdYM0h/AhuU0w3nuNX9CLPgm',
    'ST12345',
    'ادمین',
    'ادمینی',
    '00000000-0000-4000-a000-000000000001',
    '00000000-0000-4000-a000-000000000002',
    'male',
    true,
    true
) ON CONFLICT DO NOTHING;