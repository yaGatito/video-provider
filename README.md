
===================================
PROBLEMS:
    [x] validation
    [x] normalization
    [x] check if json tags are necessary
    [x] http adapter validation
    [x] service validation
    [x] file/package structure
    [x] error handling
    [] add password hashing
    [] logging
    [] unit tests
    [] swagger documentation
    [] docker support
    [] jenkins support
    [] CI/CD support

====================================

**Build application**:
```bash
go build .\src\main\cmd\main.go
```
   
**Up database**:
```bash
#go run .\src\main\cmd\migrate.go ENVIRONMENT=local
docker run -d --name=user-db -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=user-records -e MYSQL_USER=admin -e MYSQL_PASSWORD=admin -p 3306:3306 -v user-db-data:/var/lib/mysql mysql:latest
```

**Run the application**:
```bash
go run .\src\main\cmd\main.go ENVIRONMENT=local USER_SERVICE_PORT=8080
```
   
**Populate the database** (optional):
```bash
go run .\src\main\cmd\seed.go ENVIRONMENT=local
```

**Access the application**:
```bash
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"email\":\"test@example.com\", \"password\":\"asd213d2d\"}"
```

CURLs
```powershell

Invoke-RestMethod -Uri http://localhost:8080/v1/users -Method Post -Headers @{ "Content-Type" = "application/json" } -Body '{ "name": "testuser", "lastname": "testPidar", "email": "test3232@example.com", "password": "asd213d2d" }'
Invoke-RestMethod -Uri http://localhost:8080/v1/users/1 -Method Get
```

```powershell
docker exec -it user-db mysql -u root -proot user-records
```

**TODO**:
- [x] CreateUser API and FindUserById API
- [x] Add http adaptervalidation 
- [x] Add service validation
- [x] Clean




