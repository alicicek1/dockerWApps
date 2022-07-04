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

EXPOSE 8083

CMD ["/dist/main"]