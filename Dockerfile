FROM golang:1.16-alpine AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go .
RUN go build -o /proxy


FROM alpine:latest
WORKDIR /
COPY --from=build /proxy /proxy
EXPOSE 8080
CMD ["/proxy"]
