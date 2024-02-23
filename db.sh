#!/bin/sh
sudo sed -i 's/ident/md5/g' /var/lib/pgsql/data/pg_hba.conf
sudo su - postgres <<EOF
psql -c "CREATE USER babuaravind WITH PASSWORD 'root';"
psql -c "CREATE DATABASE userdb;"
psql -c "GRANT ALL PRIVILEGES ON DATABASE userdb TO babuaravind;"
psql -c "ALTER USER babuaravind WITH SUPERUSER;"
EOF