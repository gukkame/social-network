-- +migrate Up
CREATE TABLE IF NOT EXISTS Notifications (
    type VARCHAR(255) NOT NULL,
    content VARCHAR(255), 
    user_id INTEGER NOT NULL,
    recipient_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY(recipient_id) REFERENCES Users(id) ON DELETE CASCADE
    FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE Notifications;