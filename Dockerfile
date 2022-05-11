# Specify base image for Go API
FROM golang:1.17

# Specify that we need to execute any commands in directory
WORKDIR /go/src/github.com/redis_docker

# Copy everything from this project into the filesystem of the container.
COPY . .

COPY config.yaml .

# Obtain package needed to run redis commands.
RUN go get github.com/go-redis/redis

# Compile the binary EXE for our app.
RUN go build -o main .

EXPOSE 8000

# Start it
CMD [ "./main" ]

