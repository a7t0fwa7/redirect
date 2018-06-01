# build stage
FROM golang:alpine AS build-env
COPY . /go/src/github.com/polyverse/redirect
WORKDIR /go/src/github.com/polyverse/redirect
RUN GOOS=linux CGO_ENABLED=0 go build 

# final stage
FROM scratch
WORKDIR /
COPY --from=build-env /go/src/github.com/polyverse/redirect/redirect /
ENTRYPOINT ["/redirect"]

