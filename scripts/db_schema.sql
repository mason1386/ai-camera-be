-- 1. Setup Base
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Auth & Admin Module
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    phone VARCHAR(20),
    role_id UUID REFERENCES roles(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'active',
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    table_name VARCHAR(50),
    record_id VARCHAR(50),
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at DESC);

-- 3. Infrastructure Module
CREATE TABLE zones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cameras (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    zone_id UUID REFERENCES zones(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    ip_address VARCHAR(50),
    rtsp_url TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'online',
    ai_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Core Business Module
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

-- 5. Big Data Module (Partitioning)
CREATE TABLE attendance_logs (
    id UUID DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL,
    camera_id UUID,
    check_in_time TIMESTAMP NOT NULL,
    confidence_score FLOAT,
    snapshot_url TEXT,
    is_valid BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id, check_in_time),
    
    CONSTRAINT fk_logs_student FOREIGN KEY (student_id) REFERENCES students(id),
    CONSTRAINT fk_logs_camera FOREIGN KEY (camera_id) REFERENCES cameras(id)
) PARTITION BY RANGE (check_in_time);

CREATE INDEX idx_logs_student_time ON attendance_logs(student_id, check_in_time DESC);

CREATE TABLE attendance_logs_default PARTITION OF attendance_logs DEFAULT;
CREATE TABLE attendance_logs_2024_01 PARTITION OF attendance_logs FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE attendance_logs_2024_02 PARTITION OF attendance_logs FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- 6. Analytics Module
CREATE TABLE daily_attendances (
    id SERIAL PRIMARY KEY,
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID,
    date DATE NOT NULL,
    first_check_in TIMESTAMP,
    last_check_out TIMESTAMP,
    status VARCHAR(20) DEFAULT 'absent',
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, date)
);

CREATE INDEX idx_daily_class_date ON daily_attendances(class_id, date);
CREATE INDEX idx_daily_status ON daily_attendances(status);
