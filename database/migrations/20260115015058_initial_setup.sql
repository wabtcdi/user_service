-- +goose Up
-- +goose StatementBegin
-- Users table
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       first_name VARCHAR(50) NOT NULL,
                       last_name VARCHAR(50) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       phone_number VARCHAR(20),
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMPTZ
);

-- Authentication table
CREATE TABLE user_authentications (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                      user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                      password_hash VARCHAR(255) NOT NULL,
                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      deleted_at TIMESTAMPTZ
);

-- Access levels table
CREATE TABLE access_levels (
                               id SERIAL PRIMARY KEY,
                               name VARCHAR(50) UNIQUE NOT NULL,
                               description TEXT,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               deleted_at TIMESTAMPTZ
);

-- User access levels table (many-to-many relationship)
CREATE TABLE user_access_levels (
                                    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    access_level_id INTEGER NOT NULL REFERENCES access_levels(id) ON DELETE CASCADE,
                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    deleted_at TIMESTAMPTZ,
                                    PRIMARY KEY (user_id, access_level_id)
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_user_authentications_user_id ON user_authentications(user_id);
CREATE INDEX idx_user_access_levels_user_id ON user_access_levels(user_id);
CREATE INDEX idx_user_access_levels_access_level_id ON user_access_levels(access_level_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_access_levels;
DROP TABLE IF EXISTS access_levels;
DROP TABLE IF EXISTS user_authentications;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
