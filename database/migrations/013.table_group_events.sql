-- +migrate Up
CREATE TABLE IF NOT EXISTS Group_events(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL,
    happening_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY(group_id) REFERENCES User_groups(id) 
);

-- +migrate Down
DROP TABLE Group_events;