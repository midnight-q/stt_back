#!/bin/bash

ssh -t transcriber@84.201.176.202 'PGPASSWORD="fXotHxuY" pg_dump -U stt_back -h 127.0.0.1 -p 45432 -c --if-exists -f home.sql'

scp transcriber@84.201.176.202:/home/transcriber/home.sql ./home.sql

docker cp ./home.sql pg-stt_back:/home.sql

docker exec pg-stt_back sh -c "psql -U stt_back -d stt_back -f /home.sql"

rm home.sql
