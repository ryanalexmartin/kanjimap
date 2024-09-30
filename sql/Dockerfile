FROM mysql:8.0

COPY setup.sql /docker-entrypoint-initdb.d/

ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=dbname
ENV MYSQL_USER=user
ENV MYSQL_PASSWORD=password
