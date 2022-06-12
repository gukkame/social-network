-- +migrate Up
CREATE TABLE IF NOT EXISTS Group_messages(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    content VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    sender_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY(sender_id) REFERENCES Users(id),
    FOREIGN KEY(group_id) REFERENCES User_groups(id)
);

-- +migrate Down
DROP TABLE Group_messages;