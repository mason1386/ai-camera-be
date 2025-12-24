-- Up
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

-- Down
DROP TABLE IF EXISTS daily_attendances;
