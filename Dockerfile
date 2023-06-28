# Build stage
FROM golang:1.20-alpine AS build
ARG version
WORKDIR /app
COPY . .
RUN go build -ldflags "-X main.version=$version" -o wilf

# Final stage
FROM alpine:3.14
RUN addgroup -S wilf && adduser -S wilf -G wilf
USER wilf
COPY --from=build /app/wilf /usr/local/bin/wilf
ENTRYPOINT ["wilf"]
