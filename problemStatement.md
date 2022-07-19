# Create Answers
A service that exposes an API (REST) which allows people to create, update, delete and retrieve answers as key-value pairs.

An answer can be defined as:

key: string
value: string
e.g. in JSON:

{
	"key" : "name",
	"value" : "John"
}

The API should expose the following endpoints:

create answer
update answer
get answer (returns the latest answer for the given key)
delete answer
get history for given key (returns an array of events in chronological order)
An event can be defined as:

event: string
data: answer
e.g. in JSON:

{
	"event" : "create",
	"data" : {
		"key": "name",
		"value": "John"
	}
}

If a user saves the same key multiple times (using update), every answer should be saved. When retrieving an answer, it should return the latest answer.

If a user tries to create an answer that already exists - the request should fail and an adequate message or code should be returned.

If an answer doesn't exist or has been deleted, an adequate message or code should be returned.

When returning history, only mutating events (create, update, delete) should be returned. The "get" events should not be recorded.

It is possible to create a key after it has been deleted. However, it is not possible to update a deleted key. For example the following event sequences are allowed:

create → delete → create → update

create → update → delete → create → update

However, the following should not be allowed:

create → delete → update
create → create


Additional questions:

How would you support multiple users?
How would you support answers with types other than string
What are the main bottlenecks of your solution?
How would you scale the service to cope with thousands of requests?