#!/bin/bash

# goose -dir migrations postgres "${PG_DSN}" up

# This will exec the CMD from your Dockerfile
exec "$@"