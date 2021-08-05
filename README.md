# Trivia Backend

This is the backend of my personal project

## How to Run

Running this project with Docker is fairly straightforward and can be done using two commands:

```console
user:~$ docker pull tuktukmaster/trivia-backend:latest
user:~$ docker run -p 8080:8080 tuktukmaster/trivia-backed:latest
```

Alternatively you can go get all the packages then run `go run main.go`
## Why I Did the Project

I did the project for two reasons. First, I really love writing Go code for whatever reason and this was a nice opportunity to write some (fairly) complex stuff given the websockets and lobbies. Second, I wanted to have some project to showcase that I knew how to program and this seemed fairly fun to playtest.

## Technologies Used on Backend

Built using Docker and Gorilla/mux for websockets and API stuff. Language used was Go.

## Next Steps

I'm currently wanting to deploy this using Terraform on AWS but I haven't had time yet this summer. I also need to add testing and ideally have a complete CI/CD pipeline built. Additionally there's a lot of polishing that needs to be done.

# Contributors

- Berke Lunstad