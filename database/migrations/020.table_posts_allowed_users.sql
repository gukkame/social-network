-- +migrate Up
CREATE TABLE IF NOT EXISTS Posts_allowed_users(
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES Users(id),
    FOREIGN KEY(post_id) REFERENCES Posts(id) ON DELETE CASCADE
);

PRAGMA foreign_keys = ON;

-- +migrate Down
DROP TABLE Posts_allowed_users;