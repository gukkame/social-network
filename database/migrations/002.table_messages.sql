-- +migrate Up
CREATE TABLE IF NOT EXISTS Messages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  content VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  sender_id INTEGER NOT NULL,
  recipient_id INTEGER NOT NULL,
  FOREIGN KEY(sender_id) REFERENCES Users(id),
  FOREIGN KEY(recipient_id) REFERENCES Users(id)
);

-- +migrate Down
DROP TABLE Messages;