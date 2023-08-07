# IGP iw

## Running/usage
To run everything needed for the application, run this in the root of the repository

```bash
docker-compose up -d
```

This builds the api service, the notifications service and runs Postgres and Kafka.
## API
API service provides the following endpoints:
	/register
	/verify/:email/:code
	/login
	/protected

### Features
- Email verification using notification service with notification sent over Kafka. 
  After registration another email notification is sent- Welcome email
- Auth using JWT and a test endpoint /protected that checks the JWT

## Notifications
Notification service accepts messages over a Kafka topic (any MQ or Pubsub can be implemented) and in a consumer group processes these notifications.
Currently Email notification is implemented and SMS notifications mocked. Any notification type processor can be implemented.


### Features

Unit tests are not written at this time, but the code is written with tests in mind and everything should be easily testable.
