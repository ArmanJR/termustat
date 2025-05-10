-- Remove test data

-- Remove the test user
DELETE FROM users WHERE id = '00000000-0000-4000-a000-000000000003';

-- Remove the test faculty
DELETE FROM faculties WHERE id = '00000000-0000-4000-a000-000000000002';

-- Remove the test university
DELETE FROM universities WHERE id = '00000000-0000-4000-a000-000000000001';