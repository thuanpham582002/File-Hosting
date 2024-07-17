# image build
FROM golang:1.23rc1-alpine AS builder

RUN apk update \
    && apk upgrade

WORKDIR /app
COPY . .
#COPY  ../go.mod ../go.mod/go.sum ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" .


# Image chay
FROM golang:1.23rc1-alpine

RUN apk update \
    && apk upgrade

WORKDIR /app

COPY --from=builder /app/file_hosting_upload /app/
COPY --from=builder /app/templates /app/templates

# Prepare isolate environment
RUN apk add libc-dev libcap-dev git
#RUN mkdir /usr/local/etc

ENV WORKDIR /app

#COPY --from=builder /app/openapi /app/openapi

CMD ["./file_hosting_upload"]
