FROM golang:1.23.2-alpine AS build

WORKDIR /app

RUN apk update && apk add --no-cache git ca-certificates 
RUN addgroup -S regular && adduser -S regular -G regular
RUN chown regular:regular /app

COPY internal internal
COPY cmd cmd
COPY pkg pkg
COPY db db
COPY docs docs
COPY go.mod go.mod
COPY go.sum go.sum

ENV CGO_ENABLED=0
RUN go build -o /tmp/server ./cmd

FROM scratch AS deploy

COPY --from=build /tmp/server /usr/local/bin/server
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
CMD ["/usr/local/bin/server"]