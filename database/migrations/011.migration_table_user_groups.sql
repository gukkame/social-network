-- +migrate Up
CREATE TABLE IF NOT EXISTS User_groups(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    creator_id INTEGER NOT NULL,
    FOREIGN KEY(creator_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE User_groups;