# Build stage
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap .

# Runtime stage - Lambda provided.al2023 base image
FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /app/bootstrap /var/runtime/bootstrap

CMD ["bootstrap"]
