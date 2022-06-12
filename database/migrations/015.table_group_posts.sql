-- +migrate Up
CREATE TABLE IF NOT EXISTS Group_posts(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL,
    image VARCHAR(255) NULL,
    created_at DATETIME NOT NULL,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULl,
    FOREIGN KEY(group_id) REFERENCES User_groups(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE Group_posts;