# FILEPATH: /home/insom/git/kanjimap/sql/update_characters_table.py
import json
import mysql.connector

class Character:
	def __init__(self, serial, word):
		self.serial = serial
		self.word = word

# Read the JSON file
with open("chinese_characters.json") as json_file:
	data = json.load(json_file)

# Parse the JSON file into a list of Character objects
characters = []
for item in data:
	character = Character(item["serial"], item["word"])
	characters.append(character)

# Connect to the database
db = mysql.connector.connect(
	host="localhost",
	user="user",
	password="password",
	database="kanjimap"
)

# Insert each character into the database
cursor = db.cursor()
for character in characters:
	# Check if character_id exists in the database:
	if character.serial == 0:
		continue
	sql = "SELECT * FROM characters WHERE character_id = %s"
	values = (character.serial,)
	cursor.execute(sql, values)
	result = cursor.fetchone()
	if result:
		print(f"Character {character.serial} already exists")
		continue

	# If the character_id exists, but the character does not match, update the character
	if result and result[1] != character.word:
		sql = "UPDATE characters SET chinese_character = %s WHERE character_id = %s"
		values = (character.word, character.serial)
		cursor.execute(sql, values)
		print(f"Updated character {character.serial}: {character.word}")
		continue

	# If the character_id does not exist, insert the character
	sql = "INSERT INTO characters (character_id, chinese_character) VALUES (%s, %s) ON DUPLICATE KEY UPDATE chinese_character = IF(character_id = VALUES(character_id), VALUES(chinese_character), chinese_character)"
	values = (character.serial, character.word)
	cursor.execute(sql, values)
	print(f"Inserted character {character.serial}: {character.word}")

db.commit()
db.close()
