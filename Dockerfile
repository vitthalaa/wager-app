# Build-stage.
FROM golang:1.18.1-alpine AS build

ENV CGO_ENABLED=0

# Copy the project directly onto the image
ADD . /go/src/app/
WORKDIR /go/src/app/
COPY . ./

# download dependencies
RUN go mod download
# build binary
RUN go build -o /app main.go

# Packaging-stage.
FROM alpine:3.15

RUN apk --no-cache add ca-certificates

COPY --from=build /app ./app
COPY --from=build /go/src/app/.env ./.env

EXPOSE 8080 8080
ENTRYPOINT ["./app"]