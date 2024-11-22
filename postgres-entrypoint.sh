#!/bin/sh

psql -U myuser -d mydb -c "CREATE DATABASE ${POSTGRES_DB};"
psql -U myuser -d mydb -c "GRANT ALL PRIVILEGES ON DATABASE ${POSTGRES_DB} TO ${POSTGRES_USER};"

for file in $APP_HOME/schema/*.sql; do echo "running $file..."; psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -f $file; done