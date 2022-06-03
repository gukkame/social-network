-- +migrate Up
CREATE TABLE IF NOT EXISTS Notifications (
    type VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    recipient_id INTEGER NOT NULL,
    FOREIGN KEY(recipient_id) REFERENCES Users(id) ON DELETE CASCADE,
);

-- +migrate Down
DROP TABLE Notifications;