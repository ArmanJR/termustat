-- Universities Table
CREATE TABLE universities (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              name_en VARCHAR(255) NOT NULL,
                              name_fa VARCHAR(255) NOT NULL,
                              is_active BOOLEAN DEFAULT true NOT NULL,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Faculties Table
CREATE TABLE faculties (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           university_id UUID NOT NULL REFERENCES universities(id) ON DELETE CASCADE,
                           name_en VARCHAR(255) NOT NULL,
                           name_fa VARCHAR(255) NOT NULL,
                           short_code VARCHAR(10) NOT NULL,
                           is_active BOOLEAN DEFAULT true NOT NULL,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Professors Table
CREATE TABLE professors (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            university_id UUID NOT NULL REFERENCES universities(id) ON DELETE CASCADE,
                            name VARCHAR(255) NOT NULL,
                            normalized_name VARCHAR(255) NOT NULL,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Semesters Table
CREATE TABLE semesters (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           year INT NOT NULL,
                           term VARCHAR(6) NOT NULL CHECK (term IN ('spring', 'fall')),
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                           UNIQUE(year, term)
);

-- Users Table
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       student_id VARCHAR(20) NOT NULL UNIQUE,
                       first_name VARCHAR(100),
                       last_name VARCHAR(100),
                       university_id UUID NOT NULL REFERENCES universities(id),
                       faculty_id UUID NOT NULL REFERENCES faculties(id),
                       gender VARCHAR(6) NOT NULL CHECK (gender IN ('male', 'female')),
                       email_verified BOOLEAN DEFAULT false NOT NULL,
                       is_admin BOOLEAN DEFAULT false NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Courses Table
CREATE TABLE courses (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         university_id UUID NOT NULL REFERENCES universities(id) ON DELETE CASCADE,
                         faculty_id UUID NOT NULL REFERENCES faculties(id) ON DELETE CASCADE,
                         professor_id UUID NOT NULL REFERENCES professors(id),
                         semester_id UUID NOT NULL REFERENCES semesters(id),
                         code VARCHAR(50) NOT NULL,
                         name VARCHAR(255) NOT NULL,
                         weight INT NOT NULL,
                         capacity INT,
                         gender_restriction VARCHAR(6) NOT NULL CHECK (gender_restriction IN ('male', 'female', 'mixed')),
                         exam_start TIMESTAMP,
                         exam_end TIMESTAMP,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                         UNIQUE(university_id, code)
);

-- Course Times Table
CREATE TABLE course_times (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                              day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
                              start_time TIME NOT NULL,
                              end_time TIME NOT NULL,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- User Courses Table
CREATE TABLE user_courses (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                              course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                              semester_id UUID NOT NULL REFERENCES semesters(id),
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Password Resets Table
CREATE TABLE password_resets (
                                 token UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                 user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 expires_at TIMESTAMP NOT NULL,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Email Verifications Table
CREATE TABLE email_verifications (
                                     token VARCHAR(36) PRIMARY KEY,
                                     user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                     expires_at TIMESTAMP NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexes
CREATE INDEX idx_universities_is_active ON universities(is_active);
CREATE INDEX idx_faculties_university_id ON faculties(university_id);
CREATE INDEX idx_professors_university_id ON professors(university_id);
CREATE INDEX idx_users_university_id ON users(university_id);
CREATE INDEX idx_users_faculty_id ON users(faculty_id);
CREATE INDEX idx_courses_university_id ON courses(university_id);
CREATE INDEX idx_courses_faculty_id ON courses(faculty_id);
CREATE INDEX idx_courses_professor_id ON courses(professor_id);
CREATE INDEX idx_courses_semester_id ON courses(semester_id);
CREATE INDEX idx_course_times_course_id ON course_times(course_id);
CREATE INDEX idx_user_courses_user_id ON user_courses(user_id);
CREATE INDEX idx_user_courses_course_id ON user_courses(course_id);
CREATE INDEX idx_user_courses_semester_id ON user_courses(semester_id);
CREATE INDEX idx_password_resets_user_id ON password_resets(user_id);
CREATE INDEX idx_email_verifications_user_id ON email_verifications(user_id);
CREATE INDEX idx_email_verifications_expires_at ON email_verifications(expires_at);