docker build --tag prattl .;
docker run --env-file .env -p 8080:8081 -it prattl;
