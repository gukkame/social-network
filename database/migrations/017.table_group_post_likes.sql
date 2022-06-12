-- +migrate Up
CREATE TABLE IF NOT EXISTS Group_post_likes(
    type VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    user_id INTEGER NOT NULL,
    group_post_id INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY(group_post_id) REFERENCES Group_posts(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE Group_post_likes;