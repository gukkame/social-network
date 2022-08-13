-- +migrate Up
CREATE TABLE IF NOT EXISTS Group_events_users(
    status VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    group_events_id INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES Users(id),
    FOREIGN KEY(group_events_id) REFERENCES Group_events(id) ON DELETE CASCADE
);

PRAGMA foreign_keys = ON;

-- +migrate Down
DROP TABLE Group_events_users;