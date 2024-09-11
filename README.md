# Auth service test task

## API:
- [POST] /auth/signin?userID=... Get tokens by userID

<= Gets userID in query params

=> Response:
```json
{
	"accessToken": "token",
	"refreshToken": "token",
}
```

- [POST] /auth/refresh Refreshing tokens

<= Request body:
```json
{
	"accessToken": "token",
	"refreshToken": "token",
}
```

=> Response:
```json
{
	"accessToken": "token",
	"refreshToken": "token",
}
```

## Docker
- Type command `make run-docker-stdout` to run app via docker
- Type command `make run-docker-background` to run app via docker in background
- Type command `make stop-docker` to stop app via docker