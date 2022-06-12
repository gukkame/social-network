-- +migrate Up
CREATE TABLE IF NOT EXISTS Categories (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  title VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO Categories (title)
VALUES ('Go'), ('HTML5'), ('CSS'), ('JavaScript'), ('Vue.js');

-- +migrate Down
DROP TABLE Categories;