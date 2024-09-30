import json
import os

import mysql.connector

debug = False


class Character:
    def __init__(self, serial, word):
        self.serial = serial
        self.word = word


db = mysql.connector.connect(
    host=os.environ.get("MYSQL_HOST", "db"),
    user=os.environ.get("MYSQL_USER", "user"),
    password=os.environ.get("MYSQL_PASSWORD", "password"),
    database=os.environ.get("MYSQL_DATABASE", "kanjimap"),
)

cursor = db.cursor()

# Create database if it doesn't exist
cursor.execute("CREATE DATABASE IF NOT EXISTS kanjimap")
cursor.execute("USE kanjimap")

# Create tables
cursor.execute(
    """
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        token VARCHAR(255),
        email VARCHAR(255)
    )
"""
)

cursor.execute(
    """
    CREATE TABLE IF NOT EXISTS characters (
        character_id VARCHAR(255),
        chinese_character VARCHAR(255),
        PRIMARY KEY (character_id)
    )
"""
)

cursor.execute(
    """
    CREATE TABLE IF NOT EXISTS user_character_progress (
        user_id INT,
        character_id VARCHAR(255),
        learned BOOLEAN,
        PRIMARY KEY (user_id, character_id),
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (character_id) REFERENCES characters(character_id)
    )
"""
)

cursor.execute(
    """
    CREATE TABLE IF NOT EXISTS character_metadata (
        id INT AUTO_INCREMENT PRIMARY KEY,
        chinese_character VARCHAR(255) UNIQUE NOT NULL,
        frequency INT,
        cumulative_frequency FLOAT,
        pinyin VARCHAR(255),
        english TEXT
    );
"""
)

# Read the JSON file
with open("chinese_characters.json", encoding="utf-8") as json_file:
    data = json.load(json_file)

# Parse the JSON file into a list of Character objects
characters = [Character(item["serial"], item["word"]) for item in data]

# Insert or update each character in the database
for character in characters:
    sql = """
    INSERT INTO characters (character_id, chinese_character)
    VALUES (%s, %s)
    ON DUPLICATE KEY UPDATE
    chinese_character = IF(chinese_character != VALUES(chinese_character), VALUES(chinese_character),
    chinese_character)
    """
    values = (character.serial, character.word)
    try:
        cursor.execute(sql, values)
        if cursor.rowcount > 0:
            if (debug):
                print(f"Inserted/Updated character {character.serial}: {character.word}")
            else:
                pass
        else:
            if (debug):
                print(f"No changes for character {character.serial}: {character.word}")
            else:
                pass
    except mysql.connector.Error as err:
        print(f"Error inserting/updating character {character.serial}: {err}")

db.commit()
db.close()
