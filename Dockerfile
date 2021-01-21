FROM golang:1.15-alpine
WORKDIR /app
COPY .. .

RUN go build -o /docker-checker

CMD ["/docker-checker"]