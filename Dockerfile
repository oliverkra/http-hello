FROM golang:1.24.5-alpine3.22 AS build

WORKDIR /app

COPY go.* .
COPY main.* .

RUN go build -o http-hello

FROM alpine:3.22
COPY --from=build /app/http-hello /bin/http-hello
EXPOSE 8080
ENTRYPOINT [ "/bin/http-hello" ]