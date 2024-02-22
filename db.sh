#!/bin/sh
sudo su - postgres <<EOF
psql -c "CREATE USER csye6225 WITH PASSWORD 'root';"
psql -c "CREATE DATABASE userdb;"
psql -c "GRANT ALL PRIVILEGES ON DATABASE userdb TO csye6225;"
EOF
