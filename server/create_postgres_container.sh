set -a; source ./.env; set +a

docker pull postgres
docker run --name $DATABASE_DOCKER_NAME \
  -e POSTGRES_PASSWORD=$DATABASE_PASSWORD \
  -e POSTGRES_USER=$DATABASE_USER \
  -e POSTGRES_DB=$DATABASE_NAME \
  -p 5432:5432 \
  -d postgres