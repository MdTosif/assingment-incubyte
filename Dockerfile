# Stage 1: build React app into web/build
FROM node:20-alpine AS web-build
WORKDIR /app
COPY web/package.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Stage 2: compile Go binary and place UI assets in public/
FROM golang:1.24-alpine AS go-build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY cmd/ ./cmd/
COPY --from=web-build /app/build ./public

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Stage 3: minimal runtime image
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=go-build /server /server
COPY --from=go-build /src/public /public
ENV PUBLIC_DIR=/public
EXPOSE 8080
ENTRYPOINT ["/server"]
