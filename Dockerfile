FROM golang:1.17

# Set CWD inside of container
WORKDIR $GOPATH/src/github.com/TheGoldenGator/API

# Copy everything from current directory to the PWD
COPY . .

# Download all dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Expost container
EXPOSE 8080

# Run it
CMD ["API"]