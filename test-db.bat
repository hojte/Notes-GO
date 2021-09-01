docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=notes_test mysql:8.0.3
sleep 10
go test -v
pause