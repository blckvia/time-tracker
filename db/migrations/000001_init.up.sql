-- users table
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       passport_series VARCHAR(4) NOT NULL,
                       passport_number VARCHAR(6) NOT NULL,
                       name VARCHAR(100) NOT NULL,
                       surname VARCHAR(100) NOT NULL,
                       patronymic VARCHAR(100),
                       address TEXT
);

-- tasks table
CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       user_id INTEGER NOT NULL REFERENCES users(id),
                       task VARCHAR(255) NOT NULL,
                       description TEXT,
                       timer BOOLEAN DEFAULT FALSE,
                       start_time TIMESTAMP,
                       end_time TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_users_passport_series_number ON users (passport_series, passport_number);
CREATE INDEX idx_tasks_user_id ON tasks (user_id);
