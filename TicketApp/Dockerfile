FROM golang:alpine
ARG ENV

WORKDIR /build
COPY . .
RUN apk add --update alpine-sdk
RUN go mod download

WORKDIR /build/src
RUN go build -tags musl -o main .
WORKDIR /dist
RUN cp /build/src/main .

EXPOSE 8084

CMD ["/dist/main"]



