FROM golang:1.21.5
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air", "-c", ".air.toml"]

