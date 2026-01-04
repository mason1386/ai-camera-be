-- Migration: 000009_final_schema_update.sql
-- Description: Cập nhật toàn bộ schema theo thiết kế FINAL DATABASE SCHEMA

-- 1. EXTENSIONS & CONFIG
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. ENUM TYPES (Sử dụng DO block để tránh lỗi nếu type đã tồn tại)
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'identity_status') THEN
        CREATE TYPE identity_status AS ENUM ('pending', 'active', 'rejected');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'attendance_status') THEN
        CREATE TYPE attendance_status AS ENUM ('late', 'on_time', 'absent', 'early_leave');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
        CREATE TYPE event_type AS ENUM ('person', 'vehicle', 'face', 'intrusion', 'loitering', 'crowd', 'fire', 'other');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_status') THEN
        CREATE TYPE event_status AS ENUM ('new', 'processing', 'resolved', 'ignored');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'person_group') THEN
        CREATE TYPE person_group AS ENUM ('employee', 'vip', 'blacklist', 'visitor', 'other');
    END IF;
END $$;

-- 3. MODULE 1: AUTHENTICATION & PERMISSIONS
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    permissions JSONB DEFAULT '[]',
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cập nhật bảng users nếu đã tồn tại, hoặc tạo mới
CREATE TABLE IF NOT EXISTS users (
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

CREATE TABLE IF NOT EXISTS audit_logs (
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

-- 4. MODULE 2: INFRASTRUCTURE (Zones & Cameras)
CREATE TABLE IF NOT EXISTS zones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cameras (
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

CREATE TABLE IF NOT EXISTS ai_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    ai_enabled BOOLEAN DEFAULT FALSE,
    ai_types event_type[] DEFAULT '{}',
    roi_zones JSONB DEFAULT '[]', 
    active_hours JSONB DEFAULT '[]',
    sensitivity INT DEFAULT 50,
    min_confidence INT DEFAULT 60,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT uniq_camera_config UNIQUE (camera_id)
);

-- 5. MODULE 3: DATA SCOPING
CREATE TABLE IF NOT EXISTS user_camera_permissions (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, camera_id)
);

CREATE TABLE IF NOT EXISTS user_zone_permissions (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    zone_id UUID NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, zone_id)
);

-- 6. MODULE 4: IDENTITY MANAGEMENT
CREATE TABLE IF NOT EXISTS identities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    group_name person_group DEFAULT 'other',
    phone_number VARCHAR(20),
    identity_card_number VARCHAR(20),
    department VARCHAR(100),
    metadata JSONB,
    status identity_status DEFAULT 'pending',
    note TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT idx_identities_code_unique_final UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS identity_faces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    quality_score FLOAT DEFAULT 0.0,
    blur_score FLOAT DEFAULT 0.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 7. MODULE 5: EVENTS & LOGS (High Traffic - Partitioned)
-- Lưu ý: Partitioning yêu cầu các bảng mẫu và default được tạo sau
CREATE TABLE IF NOT EXISTS recognition_logs (
    id UUID DEFAULT uuid_generate_v4(), 
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    identity_id UUID REFERENCES identities(id) ON DELETE SET NULL,
    snapshot_url TEXT,
    face_crop_url TEXT,
    confidence FLOAT,
    label VARCHAR(50),
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, occurred_at)
) PARTITION BY RANGE (occurred_at);

CREATE TABLE IF NOT EXISTS ai_events (
    id UUID DEFAULT uuid_generate_v4(),
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    event_type event_type NOT NULL,
    confidence FLOAT DEFAULT 0.0,
    snapshot_url TEXT,
    metadata JSONB,
    status event_status DEFAULT 'new',
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Partitions
CREATE TABLE IF NOT EXISTS recognition_logs_default PARTITION OF recognition_logs DEFAULT;
CREATE TABLE IF NOT EXISTS ai_events_default PARTITION OF ai_events DEFAULT;

-- 8. MODULE 6: ANALYTICS (Attendance)
CREATE TABLE IF NOT EXISTS attendance_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    check_in TIMESTAMP WITH TIME ZONE,
    check_out TIMESTAMP WITH TIME ZONE,
    work_hours FLOAT DEFAULT 0,
    status attendance_status DEFAULT 'absent',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT uniq_attendance_daily_final UNIQUE (identity_id, date)
);

-- 9. INDEXES
CREATE INDEX IF NOT EXISTS idx_identities_group_final ON identities(group_name);
CREATE INDEX IF NOT EXISTS idx_identity_faces_identity_final ON identity_faces(identity_id);
CREATE INDEX IF NOT EXISTS idx_rec_logs_camera_time_final ON recognition_logs(camera_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_ai_events_type_time_final ON ai_events(event_type, created_at DESC);

-- 10. TRIGGERS
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_users_modtime ON users;
CREATE TRIGGER update_users_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

DROP TRIGGER IF EXISTS update_cameras_modtime ON cameras;
CREATE TRIGGER update_cameras_modtime BEFORE UPDATE ON cameras FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

DROP TRIGGER IF EXISTS update_identities_modtime ON identities;
CREATE TRIGGER update_identities_modtime BEFORE UPDATE ON identities FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

DROP TRIGGER IF EXISTS update_ai_configs_modtime ON ai_configs;
CREATE TRIGGER update_ai_configs_modtime BEFORE UPDATE ON ai_configs FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

DROP TRIGGER IF EXISTS update_attendance_modtime ON attendance_records;
CREATE TRIGGER update_attendance_modtime BEFORE UPDATE ON attendance_records FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

DROP TRIGGER IF EXISTS update_ai_events_modtime ON ai_events;
CREATE TRIGGER update_ai_events_modtime BEFORE UPDATE ON ai_events FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
