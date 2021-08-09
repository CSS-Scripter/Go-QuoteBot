FROM golang:1.17rc2-buster
WORKDIR app
COPY . .
RUN go mod vendor
CMD ["go", "run", "main.go"]
