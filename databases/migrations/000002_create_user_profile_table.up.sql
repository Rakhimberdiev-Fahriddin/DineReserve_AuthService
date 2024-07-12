CREATE TABLE IF NOT EXISTS user_profiles (
    user_id UUID NOT NULL REFERENCES users(id),
    username VARCHAR(100) UNIQUE,
    fullname VARCHAR(200),
    date_of_birth DATE,
    phone_number VARCHAR(20),
    address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)