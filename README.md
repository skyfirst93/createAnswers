# Create Answers
A service that exposes an API (REST) which allows people to create, update, delete and retrieve answers as key-value pairs.

## Additional questions:

#### How would you support multiple users?
#### How would you support answers with types other than string
#### What are the main bottlenecks of your solution?
#### How would you scale the service to cope with thousands of requests?

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