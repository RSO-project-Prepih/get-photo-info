# get-photo-info
This is a simple microservice that returns information about a photo and uploads it to the database. It also checks if the photo is valid meaning that it is in the correct coordinates. It use grpc to communicate with other microservices.  
## GO status
[![Go Report Card](https://goreportcard.com/badge/github.com/RSO-project-Prepih/get-photo-info)](https://goreportcard.com/report/github.com/RSO-project-Prepih/get-photo-info)

## CircleCI status CI/CD
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/RSO-project-Prepih/get-photo-info/tree/main.svg?style=svg)]

## Swagger openapi documentation
To see the swagger documentation you need to run the application and go to the following endpoint:
```
http://localhost:8080/openapi/index.html
```