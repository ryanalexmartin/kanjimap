CREATE DATABASE IF NOT EXISTS kanjimap;
GRANT ALL PRIVILEGES ON kanjimap.* TO 'user'@'%';
FLUSH PRIVILEGES;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS user_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_user_id ON user_tokens (user_id);
CREATE INDEX idx_token ON user_tokens (token);

CREATE TABLE characters (
    character_id VARCHAR(255),
    chinese_character VARCHAR(255),
    PRIMARY KEY (character_id)
);

CREATE TABLE user_character_progress (
    user_id INT,
    character_id VARCHAR(255),
    learned BOOLEAN,
    PRIMARY KEY (user_id, character_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);
