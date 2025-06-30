# To deploy the application there should be enough to have a docker

# Steps to deploy:
1. Clone a project.
2. Enter root folder (go-test-assesment) in terminal.
3. Enter command "docker-compose up --build -d"
4. After all containers were created you can use next links:
 - http://localhost:8080/swagger/index.html - for swagger docs review and test
 - http://localhost:8081/ - for pgAdmin usage (user admin@example.com, pass admin123)
5. For shutdown and deleting all created infrastructure you can use command "docker-compose down -v"

# Additionally: 
To run unit tests which are checking basic functionality you can use "go test -v ./... " before build 