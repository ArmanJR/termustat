-- Insert a test professor first
INSERT INTO professors (id, university_id, name, normalized_name)
VALUES (
           '00000000-0000-4000-a000-000000000004',
           '00000000-0000-4000-a000-000000000001', -- Test University ID
           'استادي استاد',
           'استادی استاد'
       ) ON CONFLICT DO NOTHING;

-- Insert a test semester
INSERT INTO semesters (id, year, term)
VALUES (
           '00000000-0000-4000-a000-000000000005',
           1404,
           'spring'
       ) ON CONFLICT DO NOTHING;

-- Insert a test course
INSERT INTO courses (
    id,
    university_id,
    faculty_id,
    professor_id,
    semester_id,
    code,
    name,
    weight,
    capacity,
    gender_restriction,
    exam_start,
    exam_end
)
VALUES (
           '00000000-0000-4000-a000-000000000006',
           '00000000-0000-4000-a000-000000000001', -- Test University ID
           '00000000-0000-4000-a000-000000000002', -- Test Faculty ID
           '00000000-0000-4000-a000-000000000004', -- Test Professor ID
           '00000000-0000-4000-a000-000000000005', -- Test Semester ID
           '1211003_01',
           'درس تست',
           3,
           50,
           'mixed',
           '2024-06-01 09:00:00+00',
           '2024-06-01 12:00:00+00'
       ) ON CONFLICT DO NOTHING;

-- Insert course times (two sessions per week)
INSERT INTO course_times (
    id,
    course_id,
    day_of_week,
    start_time,
    end_time
)
VALUES
    (
        '00000000-0000-4000-a000-000000000007',
        '00000000-0000-4000-a000-000000000006', -- Test Course ID
        1,
        '08:00:00',
        '09:30:00'
    ),
    (
        '00000000-0000-4000-a000-000000000008',
        '00000000-0000-4000-a000-000000000006', -- Test Course ID
        3,
        '08:00:00',
        '09:30:00'
    )
    ON CONFLICT DO NOTHING;