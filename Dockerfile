FROM quay.imanuel.dev/dockerhub/library---golang:1.16-alpine
WORKDIR /app
COPY . .

RUN go build -o /docker-checker

CMD ["/docker-checker"]