# Step 1: Modules caching
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY ./internal /zip_service/internal
COPY ./cmd /zip_service/cmd
COPY ./.env /zip_service/
COPY go.mod go.sum /zip_service/
WORKDIR /zip_service
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/main ./cmd/main/

# Step 3: Final
FROM alpine

COPY --from=builder /zip_service/.bin/main .
COPY --from=builder /zip_service/cmd/main/docs ./docs
COPY --from=builder /zip_service/.env .

ENTRYPOINT ["./main"]