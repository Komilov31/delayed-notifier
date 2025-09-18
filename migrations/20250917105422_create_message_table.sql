-- +goose Up
CREATE TABLE IF NOT EXISTS notifications(
    id SERIAL PRIMARY KEY,  
    text TEXT,
    status TEXT CHECK (status IN ('active', 'canceled', 'completed')),
    telegram_id INT NOT NULL,
    send_at BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS notifications;
