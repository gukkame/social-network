-- +migrate Up
CREATE TABLE IF NOT EXISTS Group_users(
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY(group_id) REFERENCES User_groups(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE Group_users;