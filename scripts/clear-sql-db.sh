#!/usr/bin/env bash
# Filename: clear-sql-db.sh

echo
echo CLEARING OUT LOCAL SQL DATABASE...
docker volume rm ${PROJECT_DIR_NAME}_pgdata # delete underlying volume that is storing data locally
# docker volume rm 51205-project_pgdata # originally
echo
