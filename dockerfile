FROM golang:1.16.5-alpine AS build
RUN apk update
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

WORKDIR /go/src/app

ENV CGO_ENABLED=0
ENV GOFLAGS="-insecure"

COPY go.mod ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

COPY . .
RUN go build -o /go/bin/app


FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/app /

EXPOSE 8080
ENTRYPOINT ["/app"]

