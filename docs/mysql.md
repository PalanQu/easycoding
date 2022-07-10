docker run -it --name mysql -v /data/mysql:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d -p 3306:3306  mysql:8.0
mysql -h localhost -P 3306 --protocol=tcp -u root -p123456
