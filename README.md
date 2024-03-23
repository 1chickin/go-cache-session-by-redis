# go-cache-session-by-redis
Go-based Web Server with Session Caching, API Rate Limiting and API Tracking with REDIS

## Components
- Gin (Go): Serves as the HTTP web server, handling requests and responses.
- Redis: Session caching, api rate limiting, and API call tracking with HyperLogLog for counting approximate the number of user calling api.
- Relational Database: Stores user data and session information.

## Design Highlights
- Session Management: Sessions are initially stored in both the database and Redis. Upon user login, sessions are checked/created, ensuring a single active session per user.
- Rate Limiting: For API calls to prevent abuse, using Redis to track the number of requests per user per time (5s).
- API Tracking: Utilizes Redis' HyperLogLog to efficiently estimate the number of users calling the API.
- Caching Strategy: Prioritizes Redis for session validation to enhance performance, with database lookups as a fallback mechanism.

## API Design
- /login: Authenticates users, creates a new session in DB & Redis (removing the old session if exist), and returns a session token.
- /ping: A rate-limited API that simulates processing delay, tracks calling api.
- /top: Returns the top 10 users based on the frequency of API calls.
- /count: Provides an approximate count of users who have called the /ping API, leveraging HyperLogLog.

### `/login`

- Method: POST
- Description: Authenticates users, creates a new session in the DB & Redis (removing the old session if it exists), and returns a session token.
- Request Body:
  ```json
  {
    "username": "user1",
    "password": "pass123"
  }
  ```
- Response 200 OK:
  ```json
  {
    "message": "Login successful",
    "sessionToken": "<session_token>"
  }
  ```
- Responses 401 Unauthorized:
  ```json
  {
    "error": "Invalid credentials"
  }
  ```

### `/ping`
- Method: GET
- Description: A rate-limited API that simulates a processing delay and tracks API calls.
- Headers:
- Authorization: Bearer <session_token>
- Response 200 OK: 
```json
{}
```
- Response 429 Too Many Requests:
```json
{
  "error": "Rate limit exceeded"
}
```

### `/top`
- Method: GET
- Description: Returns the top 10 users based on the frequency of API calls.
- Responses 200 OK:
```json
	{
		"topUsersCallingAPIAllTime": [
			"CallingPingAPI userID:1 called 1 times",
			"CallingPingAPI userID:3 called 4 times"
		]
	}
```
### `/count`
- Method: GET
- Description: Provides an approximate count of users who have called the /ping API, leveraging HyperLogLog.
- Response 200 OK:
```json
{
  "estimatedCount": 150
}
```

## Session Handling
- Session Validation: Prioritizes Redis for faster session validation. If a session is not found or expired in Redis, it falls back to the database check. Valid sessions found in the database but not in Redis are re-cached.
- New Session Creation: On login, any existing session for the user is removed from both Redis and the database to ensure a single active session before creating a new one with expiration time in Database and TTL in Redis.
