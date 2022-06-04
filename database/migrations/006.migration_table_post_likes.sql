-- +migrate Up
CREATE TABLE IF NOT EXISTS Post_likes (
  type VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY(post_id) REFERENCES Posts(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES Users(id)
);

-- +migrate Down
DROP TABLE Post_likes;