-- +migrate Up
CREATE TABLE IF NOT EXISTS Group_post_comment_likes(
    type VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    user_id INTEGER NOT NULL,
    group_post_comment_id INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES Users(id),
    FOREIGN KEY(group_post_comment_id) REFERENCES Group_post_comments(id) ON DELETE CASCADE
);

PRAGMA foreign_keys = ON;

-- +migrate Down
DROP TABLE Group_post_comment_likes;