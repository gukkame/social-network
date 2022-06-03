-- +migrate Up
CREATE TABLE IF NOT EXISTS Comment_likes (
  type VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  comment_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY(comment_id) REFERENCES Comments(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES Users(id)
);

-- +migrate Down
DROP TABLE Comment_likes;