#!/bin/sh
sudo sed -i 's/ident/md5/g' /var/lib/pgsql/data/pg_hba.conf
sudo su - postgres <<EOF
psql -c "CREATE USER  WITH PASSWORD 'root';"
psql -c "CREATE DATABASE ;"
psql -c "GRANT ALL PRIVILEGES ON DATABASE  TO ;"
psql -c "ALTER USER  WITH SUPERUSER;"
EOF