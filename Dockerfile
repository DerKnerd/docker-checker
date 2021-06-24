FROM quay.imanuel.dev/dockerhub/library---golang:1.16-alpine as build
WORKDIR /app
COPY . .

RUN go build -o /docker-checker .

FROM quay.imanuel.dev/dockerhub/library---alpine:latest

COPY --from=build /docker-checker /docker-checker

CMD ["/docker-checker"]