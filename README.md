# infracloud
A Golang Application to URL shortener service that will accept a URL as an argument over a REST API and return a shortened URL as a result.

## To Run the Application run the below command

Before running the application , download the dependencies using the below command.

go mod init urlShortener.

go mod tidy.

go run main.go.

Once the Golang Application started , you can use the url-shorten api's using the below URL.
http://localhost:8080/shorten    .
http://localhost:8080/metrics     .
http://localhost:8080/{shortenURL} .
http://localhost:8080/matrics/list  .

## Go throgh the Dockerfile/ Docker Hub link and pull the Docker image using the below command

docker pull maharaj2113/url-shortener:latest

