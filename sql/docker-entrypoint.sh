#!/bin/bash
set -e

# Run the Python script
python3 /docker-entrypoint-initdb.d/update_characters_table.py
echo "update_characters_table.py completed."
python3 /docker-entrypoint-initdb.d/import_frequency_data.py
echo "import_frequency_data.py completed."


# Keep the container running with MySQL in the foreground
wait $!
