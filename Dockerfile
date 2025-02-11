FROM golang:1.13
WORKDIR /app
COPY main /app/launcher/main
COPY handler /app/launcher/handler
COPY template /app/launcher/template
COPY test /app/launcher/test
COPY typed /app/launcher/typed
COPY vendor /app/launcher/vendor
COPY go.mod /app/launcher/go.mod
COPY go.sum /app/launcher/go.sum
RUN cd launcher/main && go build -mod=vendor main.go
