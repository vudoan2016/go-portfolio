FROM golang:1.16-alpine AS build

WORKDIR /go/src

RUN apk add --update gcc
RUN apk add musl-dev

COPY go.mod go.sum main.go ./
COPY analysis ./analysis
COPY financialmodelingprep ./financialmodelingprep
COPY finhub ./finhub 
COPY input ./input
COPY output ./output
COPY models ./models

# Install library dependencies
RUN go mod download

RUN go build -o portfolio .

# Build a new single layer image
FROM alpine:latest 

COPY --from=build /go/src .
COPY --from=build /go/src/output ./output
COPY --from=build /go/src/input ./input
RUN touch inpfo.log

# Expose port 8080 to the outside world.
EXPOSE 8080

ENTRYPOINT ["./portfolio"]