FROM nekkkkitch/docker AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux go build -a -o register ./services/register/cmd/app/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/register .
COPY ./services/register/cfg.yml .
CMD ["./register"]