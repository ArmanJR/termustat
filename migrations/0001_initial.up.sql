CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE universities (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              name VARCHAR(255) NOT NULL,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE faculties (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           university_id UUID REFERENCES universities(id) ON DELETE CASCADE,
                           name VARCHAR(255) NOT NULL,
                           short_code VARCHAR(50) NOT NULL,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       student_id VARCHAR(50) NOT NULL,
                       first_name VARCHAR(100),
                       last_name VARCHAR(100),
                       university_id UUID REFERENCES universities(id),
                       faculty_id UUID REFERENCES faculties(id),
                       gender VARCHAR(6) CHECK (gender IN ('male', 'female')),
                       email_verified BOOLEAN DEFAULT false,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE courses (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         university_id UUID REFERENCES universities(id) ON DELETE CASCADE,
                         faculty_id UUID REFERENCES faculties(id) ON DELETE CASCADE,
                         code VARCHAR(50) NOT NULL,
                         name VARCHAR(255) NOT NULL,
                         weight INT NOT NULL,
                         capacity INT,
                         gender_restriction VARCHAR(6) CHECK (gender_restriction IN ('male', 'female', 'mixed')),
                         professor VARCHAR(255),
                         exam_start TIMESTAMP,
                         exam_end TIMESTAMP,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         UNIQUE(university_id, code)
);

CREATE TABLE course_schedules (
                                  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
                                  day_of_week INT CHECK (day_of_week BETWEEN 0 AND 6),
                                  start_time TIME NOT NULL,
                                  end_time TIME NOT NULL
);

CREATE TABLE user_courses (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                              course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
                              semester VARCHAR(50) NOT NULL,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE password_resets (
                                 token UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                 user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                 expires_at TIMESTAMP NOT NULL
);