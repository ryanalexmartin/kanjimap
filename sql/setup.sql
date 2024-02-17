CREATE DATABASE IF NOT EXISTS kanjimap;
USE kanjimap;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255),
    email VARCHAR(255) NOT NULL,
);

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
