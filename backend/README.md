# Environment Variables
As an attempt to look pro, this backend project uses environment variables
to save your secrets, here are all the required ones

## SERVER_PORT
It is used to define the port you want to use on the server

## CONNECTION_STRING
It specifies the connection to the database using the postgres protocol

## AUTH_KEY
It is the key you use to sign your Json Web Tokens, it should be a secret

## ALLOWED_HOST
It is the frontend host you want to allow to perform requests, it usually will be
http://localhost:4200 (For Angular)

# .env file example


``
SERVER_PORT=":3000"
CONNECTION_STRING="postgresql://tasks:tasks@localhost?sslmode=disable"
AUTH_KEY="super-secret-auth-key"
ALLOWED_HOST="http://localhost:4200"
``
