/* DATABASE TEMPLATE FOR RECREATING THE DATABASE IN CASE OF A MISSING database.db FILE */

/* FOREIGN KEY CONSTRAIN ON: NEEDED FOR 'ON DELETE CASCADE' TO WORK */
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Users (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  age VARCHAR(255) NOT NULL,
  gender VARCHAR(255) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS Messages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  content VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  sender_id INTEGER NOT NULL,
  recipient_id INTEGER NOT NULL,
  FOREIGN KEY(sender_id) REFERENCES Users(id),
  FOREIGN KEY(recipient_id) REFERENCES Users(id)
);

CREATE TABLE IF NOT EXISTS Sessions (
  token VARCHAR(255) NOT NULL PRIMARY KEY,
  expiry_date DATETIME NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Categories (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  title VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO Categories (title)
VALUES ('Go'), ('HTML5'), ('CSS'), ('JavaScript'), ('Vue.js');



CREATE TABLE IF NOT EXISTS Posts (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  title VARCHAR(255) NOT NULL,
  content VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  category_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY(category_id) REFERENCES Categories(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES Users(id)
);

CREATE TABLE IF NOT EXISTS Post_likes (
  type VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY(post_id) REFERENCES Posts(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES Users(id)
);

CREATE TABLE IF NOT EXISTS Comments (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  content VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  FOREIGN KEY(user_id) REFERENCES Users(id),
  FOREIGN KEY(post_id) REFERENCES Posts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Comment_likes (
  type VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  comment_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY(comment_id) REFERENCES Comments(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES Users(id)
);

/* FOREIGN KEY CONSTRAIN ON: NEEDED FOR 'ON DELETE CASCADE' TO WORK */
PRAGMA foreign_keys = ON;



