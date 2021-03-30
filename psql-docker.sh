#!/bin/bash
docker run --rm --name pg-stt_back -e POSTGRES_DB=stt_back -e POSTGRES_USER=stt_back -e POSTGRES_PASSWORD=fXotHxuY -d -p 45432:5432 -v "$(pwd)/postgres:/var/lib/postgresql/data" postgres
