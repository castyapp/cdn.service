FROM golang:1.15-alpine AS builder

LABEL maintainer="Alireza Josheghani <josheghani.dev@gmail.com>"

# Creating work directory
WORKDIR /app

# Adding project to work directory
ADD . /app

# build project
RUN go build -o server .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server /app

EXPOSE 5555

ENTRYPOINT ["/app/server"]
CMD ["--port", "5555"]