-- +migrate Up
CREATE TABLE IF NOT EXISTS Notifications (
    type VARCHAR(255) NOT NULL,
    content_id INTEGER, 
    recipient_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY(recipient_id) REFERENCES Users(id) ON DELETE CASCADE
    FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE Notifications;