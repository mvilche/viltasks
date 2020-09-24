docker-compose down
docker-compose up -d
docker-compose exec viltasks bash -c "tar zcvf /tmp/app.tar.gz /app"
docker cp viltasks_viltasks_1:/tmp/app.tar.gz .