FROM golang:1.13
WORKDIR /app
COPY main /app/launcher/main
COPY rest /app/launcher/rest
COPY template /app/launcher/template
COPY test /app/launcher/test
COPY typed /app/launcher/typed
COPY vendor /app/launcher/vendor
COPY go.mod /app/launcher/go.mod
COPY go.sum /app/launcher/go.sum
COPY config /app/.kube
RUN cd launcher/main && go build main.go
RUN cd launcher/test && go build test.go
