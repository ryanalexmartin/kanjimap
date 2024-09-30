import csv

import mysql.connector
from mysql.connector import Error

debug = False


def import_frequency_data(file_path, db_config):
    try:
        connection = mysql.connector.connect(**db_config)
        cursor = connection.cursor()

        # First, let's get all existing characters from the characters table
        cursor.execute("SELECT chinese_character FROM characters")
        existing_characters = set(row[0] for row in cursor.fetchall())

        # Prepare the update statement
        update_query = """
        UPDATE character_metadata
        SET frequency = %s, cumulative_frequency = %s, pinyin = %s, english = %s
        WHERE chinese_character = %s
        """

        # Prepare the insert statement (we'll use this only for existing characters not yet in character_metadata)
        insert_query = """
        INSERT INTO character_metadata
        (chinese_character, frequency, cumulative_frequency, pinyin, english)
        VALUES (%s, %s, %s, %s, %s)
        """

        updated_count = 0
        skipped_count = 0

        with open(file_path, "r", encoding="utf-8") as file:
            tsv_reader = csv.reader(file, delimiter="\t")
            next(tsv_reader)  # Skip the header row

            for row in tsv_reader:
                serial, character, frequency, cumulative_frequency, pinyin, english = (
                    row
                )

                if character in existing_characters:
                    # Check if the character already exists in character_metadata
                    cursor.execute(
                        "SELECT 1 FROM character_metadata WHERE chinese_character = %s",
                        (character,),
                    )
                    if cursor.fetchone():
                        # Update existing record
                        cursor.execute(
                            update_query,
                            (
                                int(frequency),
                                float(cumulative_frequency),
                                pinyin,
                                english,
                                character,
                            ),
                        )
                    else:
                        # Insert new record for existing character
                        cursor.execute(
                            insert_query,
                            (
                                character,
                                int(frequency),
                                float(cumulative_frequency),
                                pinyin,
                                english,
                            ),
                        )
                    updated_count += 1
                else:
                    skipped_count += 1
                    if debug:
                        print(f"Character '{character}' not found in characters table. Skipping.")

        connection.commit()
        if debug:
            print(f"Data import completed. Updated {updated_count} characters. Skipped {skipped_count} characters.")

        # Verify the import
        cursor.execute("SELECT COUNT(*) FROM character_metadata")
        count = cursor.fetchone()[0]
        if debug:
            print(f"Total records in character_metadata: {count}")

        if debug:
            cursor.execute(
                "SELECT chinese_character, frequency FROM character_metadata ORDER BY frequency DESC LIMIT 5"
            )
            print("Top 5 most frequent characters:")
            for row in cursor.fetchall():
                print(row)

    except Error as e:
        print(f"Error: {e}")

    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


# Database configuration
db_config = {
    "host": "db",
    "user": "user",
    "password": "password",
    "database": "kanjimap",
}

# Path to your TSV file
file_path = "character_frequency.tsv"

# Run the import function
import_frequency_data(file_path, db_config)
