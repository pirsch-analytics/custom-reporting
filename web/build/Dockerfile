FROM golang:1.21 AS build
ARG TARGET=build
RUN apt-get update && \
    apt-get upgrade -y
WORKDIR /go/src/custom-report
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/custom-report/cmd/main.go && \
	mkdir /app && \
	mv /go/src/custom-report/main /app/server

RUN mv /go/src/custom-report/static /app/static

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=build /app /app
WORKDIR /app

RUN addgroup -S appuser && \
    adduser -S -G appuser -u 1000 appuser && \
    chown -R appuser:appuser /app
USER appuser

EXPOSE 8080
ENTRYPOINT ["/app/server"]
