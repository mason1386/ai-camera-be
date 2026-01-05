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
-- =================================================================
-- MIGRATION: 000007_init_face_recognition_schema.sql
-- AUTHOR: CameraAI Architect
-- DESCRIPTION: Schema cho quản lý Face, Lịch sử nhận diện (Partitioned) và Điểm danh
-- NOTE: Yêu cầu DB đã có bảng 'users' (từ migration 000001) và 'cameras' (từ migration 000004)
-- =================================================================

-- 1. Kích hoạt Extension (Cần thiết cho UUID và Vector Search sau này)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- CREATE EXTENSION IF NOT EXISTS vector; -- Uncomment nếu server đã cài pgvector

-- 2. Định nghĩa các Enum Types
CREATE TYPE identity_status AS ENUM ('pending', 'active', 'rejected');
CREATE TYPE attendance_status AS ENUM ('late', 'on_time', 'absent', 'early_leave');

-- =================================================================
-- TABLE: IDENTITIES (Hồ sơ khuôn mặt)
-- Liên kết chặt chẽ với bảng USERS
-- =================================================================
CREATE TABLE identities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Thông tin định danh cơ bản
    code VARCHAR(50) NOT NULL,              -- Mã nhân viên/học sinh/khách (Unique)
    full_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    identity_card_number VARCHAR(20),       -- CCCD/CMND
    
    -- Dữ liệu khuôn mặt & AI
    face_image_url TEXT NOT NULL,           -- URL ảnh gốc (MinIO/S3)
    -- face_embedding VECTOR(512),          -- Uncomment nếu dùng pgvector
    
    -- Phân loại & Trạng thái
    type VARCHAR(50) NOT NULL,              -- Config: STAFF, STUDENT, VIP...
    status identity_status DEFAULT 'pending',
    note TEXT,

    -- [LIÊN KẾT 1: AUDIT] - Ai là người tạo/duyệt hồ sơ này?
    created_by UUID REFERENCES users(id) ON DELETE SET NULL, 
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- [LIÊN KẾT 2: LOGIN] - Hồ sơ này gắn với tài khoản đăng nhập nào?
    -- (Để nhân viên có thể login vào App xem chấm công của chính mình)
    user_account_id UUID UNIQUE REFERENCES users(id) ON DELETE SET NULL,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT idx_identities_code_unique UNIQUE (code)
);

-- Indexes tối ưu tìm kiếm
CREATE INDEX idx_identities_type ON identities(type);
CREATE INDEX idx_identities_status ON identities(status);
CREATE INDEX idx_identities_user_acc ON identities(user_account_id); -- Tìm nhanh identity từ user_id đang login

-- =================================================================
-- TABLE: RECOGNITION_LOGS (Lịch sử nhận diện - High Traffic)
-- Chiến lược: Partitioning theo thời gian (Tháng)
-- =================================================================
CREATE TABLE recognition_logs (
    id UUID DEFAULT uuid_generate_v4(), 
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    
    -- Có thể Null nếu là người lạ (Stranger)
    identity_id UUID REFERENCES identities(id) ON DELETE SET NULL,
    
    -- Hình ảnh sự kiện
    snapshot_url TEXT,         -- Ảnh toàn cảnh
    face_crop_url TEXT,        -- Ảnh cắt khuôn mặt
    
    -- AI Metadata
    confidence FLOAT,          -- Độ chính xác (0.0 - 1.0)
    label VARCHAR(50),         -- Label tại thời điểm nhận diện (Stranger, Staff...)
    
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Partition Key phải nằm trong Primary Key
    PRIMARY KEY (id, occurred_at)
) PARTITION BY RANGE (occurred_at);

-- Tạo Partition mặc định (Chứa dữ liệu lỗi thời gian hoặc chưa phân vùng)
CREATE TABLE recognition_logs_default PARTITION OF recognition_logs DEFAULT;

-- Tạo Partition mẫu cho năm 2024, 2025 (Bạn nên viết cronjob tạo tự động)
CREATE TABLE recognition_logs_2024_12 PARTITION OF recognition_logs
    FOR VALUES FROM ('2024-12-01') TO ('2025-01-01');
CREATE TABLE recognition_logs_2025_01 PARTITION OF recognition_logs
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

-- Indexes Local cho từng Partition (Hiệu năng cao)
CREATE INDEX idx_rec_logs_camera_time ON recognition_logs(camera_id, occurred_at DESC);
CREATE INDEX idx_rec_logs_identity_time ON recognition_logs(identity_id, occurred_at DESC);

-- =================================================================
-- TABLE: ATTENDANCE_RECORDS (Dữ liệu điểm danh tổng hợp)
-- =================================================================
CREATE TABLE attendance_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    
    date DATE NOT NULL,
    
    check_in TIMESTAMP WITH TIME ZONE,  -- Lần xuất hiện đầu tiên
    check_out TIMESTAMP WITH TIME ZONE, -- Lần xuất hiện cuối cùng
    
    work_hours FLOAT DEFAULT 0,         -- Số giờ làm việc tính toán được
    status attendance_status DEFAULT 'absent',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Mỗi người chỉ có 1 dòng điểm danh mỗi ngày
    CONSTRAINT uniq_attendance_daily UNIQUE (identity_id, date)
);

CREATE INDEX idx_attendance_date_status ON attendance_records(date, status);
CREATE INDEX idx_attendance_identity_month ON attendance_records(identity_id, date);

-- =================================================================
-- TRIGGERS (Tự động cập nhật updated_at)
-- =================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_identities_modtime BEFORE UPDATE ON identities FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_attendance_modtime BEFORE UPDATE ON attendance_records FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
ALTER TABLE identities ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
