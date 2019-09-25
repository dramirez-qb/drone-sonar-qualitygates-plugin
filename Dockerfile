# Start by building the application.
FROM golang:1.13-alpine as build
WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/drone-sonar-qualitygate

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /go/bin/drone-sonar-qualitygate /
CMD ["/drone-sonar-qualitygate"]
