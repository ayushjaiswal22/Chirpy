# Chirpy

This is a learning project for golang

## Installation

`go install github.com/ayushjaiswal22/chirpy@latest`

## Configuration

### Create a .env file with below entries:

```
PLATFORM="dev"
SECRET_KEY="au98AyorIK073Kd1uwrtkwssx5gSGIVYKUvzSLgX7qr6soCqpWYMUDNYO4mVC7kAAO+mMSuiQz2jGivuKvYv0Q=="
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
```

## Run Chirpy

`go run .`

### Now you can make API calls as below

`curl -X <METHOD> http://localhost:8080/<endpoint> -d '<request_body>'`

### Chirpy supports the following methods:

| Function                  | Method    | Endpoint                        | Request JSON                                                                    | Comments |
|---------------------------|-----------|---------------------------------|---------------------------------------------------------------------------------|----------|
| Create User               | POST      | /api/users                      | `{"email":"ab@example.com", "password":"pass123"}`                              |          |
| Update User               | POST      | /api/users                      | `{"email":"ab@example.com", "password":"pass123"}`                              |          |
| Login                     | POST      | /api/login                      | `{"email":"ab@example.com", "password":"pass123", "expires_in_seconds": 100}`   |          |

#### We use JWT as access tokens for Authentication and refresh tokens in order to refresh/revoke an access token

##### When you login we create an access token 

| Function                    | Method    | Endpoint                        | Header                                          | Request JSON                                         | Comments                |
|-----------------------------|-----------|---------------------------------|-------------------------------------------------|------------------------------------------------------|-------------------------|
| Refresh an access token     | POST      | /api/refresh                    | "Authorization: Bearer your_refresh_token_here" |                                                      |                         |
| Revoke an access token      | POST      | /api/revoke                     | "Authorization: Bearer your_refresh_token_here" |                                                      |                         |
| Post Chirp                  | POST      | /api/chirps                     | "Authorization: Bearer your_access_token_here"  | `{"user_id":"123-abc-456", "body":"Hey There! :-)"}` |                         |
| Get All Chirps              | GET       | /api/chirps                     |                                                 |                                                      |                         |
| Get a Chirp by id           | GET       | /api/chirps/{chirp_id}          |                                                 |                                                      |                         |
| Get All Chirps of a user    | GET       | /api/chirps?author_id=<user_id> |                                                 |                                                      |                         |
| Delete a Chirp by id        | DELETE    | /api/chirps/{chirp_id}          | "Authorization: Bearer your_access_token_here"  |                                                      |                         |
| Webhook by Polka            | POST      | /api/polka/webhooks             | "Authorization: ApiKey Polka_API_Key_here"      |                                                      | To upgrade to Chirpy Red|


