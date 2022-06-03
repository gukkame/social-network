-- +migrate Up
CREATE TABLE IF NOT EXISTS Followers (
status VARCHAR(255) NOT NULL,
created_at DATETIME NOT NULL,
follower_id INTEGER NOT NULL,
recipient_id INTEGER NOT NULL,
FOREIGN KEY(follower_id) REFERENCES Users(id) ON DELETE CASCADE,
FOREIGN KEY(recipient_id) REFERENCES Users(id) ON DELETE CASCADE
)

-- +migrate Down
DROP TABLE Followers;