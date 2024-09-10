import csv
import mysql.connector
from mysql.connector import Error

def import_frequency_data(file_path, db_config):
    try:
        connection = mysql.connector.connect(**db_config)
        cursor = connection.cursor()

        # Prepare the insert statements
        insert_character = """
        INSERT INTO characters (character_id, chinese_character)
        VALUES (%s, %s)
        ON DUPLICATE KEY UPDATE chinese_character = VALUES(chinese_character)
        """
        
        insert_metadata = """
        INSERT INTO character_metadata 
        (character_id, frequency, cumulative_frequency, pinyin, english)
        VALUES (%s, %s, %s, %s, %s)
        ON DUPLICATE KEY UPDATE
        frequency = GREATEST(frequency, VALUES(frequency)),
        cumulative_frequency = GREATEST(cumulative_frequency, VALUES(cumulative_frequency)),
        pinyin = VALUES(pinyin),
        english = VALUES(english)
        """

        # Fetch all characters and their IDs from the characters table
        cursor.execute("SELECT character_id, chinese_character FROM characters")
        char_id_map = {row[1]: row[0] for row in cursor.fetchall()}

        print(f"Initial characters in database: {len(char_id_map)}")

        processed = 0
        added = 0
        updated = 0
        char_frequency_map = {}

        # First pass: collect all data and keep the highest frequency for each character
        with open(file_path, 'r', encoding='utf-8') as file:
            tsv_reader = csv.reader(file, delimiter='\t')
            next(tsv_reader)  # Skip the header row
            
            for row in tsv_reader:
                serial, character, frequency, cumulative_frequency, pinyin, english = row
                frequency = int(frequency)
                cumulative_frequency = float(cumulative_frequency)
                
                if character not in char_frequency_map or frequency > char_frequency_map[character][1]:
                    char_frequency_map[character] = (serial, frequency, cumulative_frequency, pinyin, english)

        # Second pass: process the collected data
        for character, (serial, frequency, cumulative_frequency, pinyin, english) in char_frequency_map.items():
            if character not in char_id_map:
                # Add the character to the characters table
                new_id = f"A{serial.zfill(5)}"
                cursor.execute(insert_character, (new_id, character))
                char_id_map[character] = new_id
                added += 1
                print(f"Added character {character} with ID {new_id}")
            else:
                updated += 1

            character_id = char_id_map[character]
            
            # Insert or update character_metadata
            cursor.execute(insert_metadata, (
                character_id,
                frequency,
                cumulative_frequency,
                pinyin,
                english
            ))
            processed += 1

            if processed % 1000 == 0:
                print(f"Processed: {processed}, Added: {added}, Updated: {updated}")

        connection.commit()
        print("Data import completed successfully.")
        print(f"Total processed: {processed}")
        print(f"Total added to characters table: {added}")
        print(f"Total updated in character_metadata: {updated}")

        # Verify the import
        cursor.execute("SELECT COUNT(*) FROM characters")
        char_count = cursor.fetchone()[0]
        cursor.execute("SELECT COUNT(*) FROM character_metadata")
        metadata_count = cursor.fetchone()[0]
        print(f"Total records in characters table: {char_count}")
        print(f"Total records in character_metadata: {metadata_count}")

    except Error as e:
        print(f"Error: {e}")

    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()

# Database configuration
db_config = {
    'host': 'localhost',
    'user': 'user',
    'password': 'password',
    'database': 'kanjimap'
}

# Path to your TSV file
file_path = 'character_frequency.tsv'

# Run the import function
import_frequency_data(file_path, db_config)
