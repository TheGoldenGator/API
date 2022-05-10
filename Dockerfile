FROM golang:1.17

# Set CWD inside of container
WORKDIR $GOPATH/src/github.com/TheGoldenGator/API

COPY go.mod ./
COPY go.sum ./
RUN go mod Download

COPY *.go ./

RUN go build -o /API

EXPOSE 8080

CMD ["api"]