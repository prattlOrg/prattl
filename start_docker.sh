docker build --tag prattl .;
docker run --env-file .env -p 8080:8080 -it prattl;
