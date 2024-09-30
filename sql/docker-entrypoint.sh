#!/bin/bash
set -e

# Start MySQL in the background
docker-entrypoint.sh mysqld &

# Wait for MySQL to be ready
# until mysqladmin ping -h"localhost" --silent; do
#     echo 'Waiting for MySQL to be ready...'
#     sleep 2
# done
#
# echo "MySQL is ready. Running update_characters_table.py..."

# Run the Python script
python3 /docker-entrypoint-initdb.d/update_characters_table.py
echo "update_characters_table.py completed."
python3 /docker-entrypoint-initdb.d/import_frequency_data.py
echo "import_frequency_data.py completed."


# Keep the container running with MySQL in the foreground
wait $!
