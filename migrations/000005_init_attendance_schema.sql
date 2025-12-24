-- Up
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

-- Down
DROP TABLE IF EXISTS attendance_logs;
