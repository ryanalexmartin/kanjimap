import json
import os

import mysql.connector
from mysql.connector import Error

debug = False


class Character:
    def __init__(self, serial, word):
        self.serial = serial
        self.word = word


def create_connection():
    try:
        connection = mysql.connector.connect(
            host=os.environ.get("MYSQL_HOST", "db"),
            user=os.environ.get("MYSQL_USER", "user"),
            password=os.environ.get("MYSQL_PASSWORD", "password"),
            database=os.environ.get("MYSQL_DATABASE", "kanjimap"),
        )
        return connection
    except Error as e:
        print(f"Error connecting to MySQL: {e}")
        return None


def execute_query(connection, query, params=None):
    cursor = connection.cursor()
    try:
        if params:
            cursor.execute(query, params)
        else:
            cursor.execute(query)
        connection.commit()
    except Error as e:
        print(f"Error executing query: {e}")
    finally:
        cursor.close()


def create_tables(connection):
    tables = [
        """
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL,
            token VARCHAR(255),
            email VARCHAR(255)
        )
        """,
        """
        CREATE TABLE IF NOT EXISTS characters (
            character_id VARCHAR(255),
            chinese_character VARCHAR(255),
            PRIMARY KEY (character_id)
        )
        """,
        """
        CREATE TABLE IF NOT EXISTS user_character_progress (
            user_id INT,
            character_id VARCHAR(255),
            learned BOOLEAN,
            PRIMARY KEY (user_id, character_id),
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (character_id) REFERENCES characters(character_id)
        )
        """,
        """
        CREATE TABLE IF NOT EXISTS character_metadata (
            id INT AUTO_INCREMENT PRIMARY KEY,
            chinese_character VARCHAR(255) UNIQUE NOT NULL,
            frequency INT,
            cumulative_frequency FLOAT,
            pinyin VARCHAR(255),
            english TEXT
        )
        """,
        """
        CREATE TABLE IF NOT EXISTS user_tokens (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            token VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
        """,
    ]

    for table in tables:
        execute_query(connection, table)

    # Check and create indexes if they don't exist
    indexes = [
        ("idx_user_id", "user_tokens", "user_id"),
        ("idx_token", "user_tokens", "token"),
    ]

    for index_name, table_name, column_name in indexes:
        if not index_exists(connection, table_name, index_name):
            create_index(connection, table_name, index_name, column_name)


def index_exists(connection, table_name, index_name):
    cursor = connection.cursor()
    try:
        cursor.execute(f"SHOW INDEX FROM {table_name} WHERE Key_name = '{index_name}'")
        return cursor.fetchone() is not None
    except Error as e:
        print(f"Error checking if index exists: {e}")
        return False
    finally:
        cursor.close()


def create_index(connection, table_name, index_name, column_name):
    try:
        execute_query(
            connection, f"CREATE INDEX {index_name} ON {table_name} ({column_name})"
        )
        print(f"Created index {index_name} on {table_name}.{column_name}")
    except Error as e:
        print(f"Error creating index {index_name}: {e}")


def insert_or_update_character(connection, character):
    sql = """
    INSERT INTO characters (character_id, chinese_character)
    VALUES (%s, %s)
    ON DUPLICATE KEY UPDATE
    chinese_character = IF(chinese_character != VALUES(chinese_character), VALUES(chinese_character), chinese_character)
    """
    values = (character.serial, character.word)
    try:
        cursor = connection.cursor()
        cursor.execute(sql, values)
        connection.commit()
        if cursor.rowcount > 0:
            if debug:
                print(
                    f"Inserted/Updated character {character.serial}: {character.word}"
                )
        else:
            if debug:
                print(f"No changes for character {character.serial}: {character.word}")
    except Error as err:
        print(f"Error inserting/updating character {character.serial}: {err}")
    finally:
        cursor.close()


def main():
    connection = create_connection()
    if connection is None:
        return

    create_tables(connection)

    with open("chinese_characters.json", encoding="utf-8") as json_file:
        data = json.load(json_file)

    characters = [Character(item["serial"], item["word"]) for item in data]

    for character in characters:
        insert_or_update_character(connection, character)

    connection.close()
    print("Database initialization completed successfully.")


if __name__ == "__main__":
    main()
