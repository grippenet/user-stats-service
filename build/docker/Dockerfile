##########################
# Build stage
##########################

ARG DIR=/go/src/github.com/grippenet/user-stats-service

FROM golang:1.17-alpine as builder
ARG DIR # Populate global ARG for this build
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN mkdir -p ${DIR}
ENV GO111MODULE=on
ADD . ${DIR}
WORKDIR ${DIR}
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR $DIR/cmd/user-stats-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

##########################
# Image 
##########################
FROM scratch
ARG DIR
COPY --from=builder $DIR/cmd/user-stats-service/app /app/
WORKDIR /app
ENV GIN_MODE=release
ENV PORT=3252
EXPOSE 3252:3252
CMD ["./app"]
