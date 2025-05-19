-- Remove the test course times
DELETE FROM course_times
WHERE id IN (
             '00000000-0000-4000-a000-000000000007',
             '00000000-0000-4000-a000-000000000008'
    );

-- Remove the test course
DELETE FROM courses
WHERE id = '00000000-0000-4000-a000-000000000006';

-- Remove the test semester
DELETE FROM semesters
WHERE id = '00000000-0000-4000-a000-000000000005';

-- Remove the test professor
DELETE FROM professors
WHERE id = '00000000-0000-4000-a000-000000000004';