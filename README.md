# Create Answers
A service that exposes an API (REST) which allows people to create, update, delete and retrieve answers as key-value pairs.

## Additional questions:

#### How would you support multiple users?
1. Using horizontal scaling 
2. Using caching 
#### How would you support answers with types other than string?
Using interface{} as type for the value data member
example Value interface{} `json:"value"`
#### What are the main bottlenecks of your solution?
1. Use a better performance providing Web Frameworks
2. Excute multiple write queeries together 
3. To work with multiple CPU streams

#### How would you scale the service to cope with thousands of requests?
1. Using horizontal scaling 
2. Using caching 
 

## Procedure to test
1. Start the Redis server 
```
	docker-compose up -d 
```

2. Start the Application 
```
	go run main.go
```

3. Test the api using the answers.postman_collection file
