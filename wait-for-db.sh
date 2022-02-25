#!/bin/sh

set -e

cmd="$@"

until PGPASSWORD=$POSTGRES_PASSWORD psql --host=$POSTGRES_HOST --username=$POSTGRES_USER -w &>/dev/null
do
  echo "Waiting for PostgreSQL..."
  sleep 1
done

>&2 echo "PostgreSQL ✔"
exec $cmd
