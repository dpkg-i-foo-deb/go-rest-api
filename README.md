# Go REST API Attempt

This is a silly attempt to create a REST API using Go, as well as
using Json Web Tokens to handle authenticated requests and generating
refresh tokens so the user doesn't have to log in all the time

As a little extra, this API communicates to a database and gets data
using raw SQL statements

I must clarify I have 0 experience when it comes to software, so I'm sure
this project doesn't follow any standard or good practice. However, my goal
is to learn a little bit more about this language and to create a useful
base for my personal projects

# What do I need to run this project
## Go
## A PostgreSQL database server
## Something to test the endpoints

That is literally all, this language is awesome

# Database ER model (IT IS NOT FULLY IMPLEMENTED)
![alt text](/database/model.png)

# How do I run this project?

- Create the database using the provided script
- Go to the "backend" directory run on the terminal...
`go run .`

That's all, the server should be running now!

# Some extras

I read it is nice to have some stuff ia .env file, so you need to create it
- Create a .env file inside the "backend" directory
- Add the required variables to it, so far you need...

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