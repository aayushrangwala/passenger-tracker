FROM golang:alpine AS build
WORKDIR /go/src/passenger-tracker
COPY ./ ./
ARG GIT_COMMIT
ARG GOPROXY

RUN go build -ldflags main.go -o /go/bin/passenger-tracker

FROM gcr.io/distroless/base-debian11
COPY --from=build /go/bin/passenger-tracker /usr/local/bin/passenger-tracker

ENTRYPOINT [ "/usr/local/bin/passenger-tracker" ]