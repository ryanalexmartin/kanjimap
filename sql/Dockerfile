# Use the official image as a parent image
FROM mysql:8.0

# Set the working directory in the container
WORKDIR /docker-entrypoint-initdb.d

# Add the database setup script
COPY setup.sql .

# Set environment variables
ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=dbname
ENV MYSQL_USER=user
ENV MYSQL_PASSWORD=password

# Expose port 3306
EXPOSE 3306

COPY update_characters_table.py .
RUN apt-get update && apt-get install -y python3 python3-pip
COPY requirements.txt . 
RUN pip3 install -r requirements.txt

# Run the Python script
CMD ["python3", "update_characters_table.py"]