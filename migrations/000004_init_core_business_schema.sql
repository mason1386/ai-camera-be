-- Up
CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    grade INT NOT NULL,
    academic_year VARCHAR(20),
    homeroom_teacher_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_code VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    date_of_birth DATE,
    gender VARCHAR(10),
    avatar_url TEXT,
    parent_phone VARCHAR(20),
    face_registered BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'studying',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_student_code ON students(student_code);
CREATE INDEX idx_student_class ON students(class_id);

-- Down
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS classes;
