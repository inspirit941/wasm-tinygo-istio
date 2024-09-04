# Use a Docker image that includes TinyGo
FROM tinygo/tinygo:latest AS builder

# Create a directory for the project
WORKDIR /project
# tinygo 빌드하려면 root 권한으로 실행해야 한다. root권한이 아니면 wasm binary 생성이 안됨
USER root

# Copy the Go source file(s) into the container
COPY . .

# Build the WASM file from the Go source code using TinyGo
# The flags indicate that we don't want a scheduler and are targeting WASI
RUN tinygo build -o /project/plugin.wasm -scheduler=none -target=wasi /project/main.go

FROM scratch
# Copy the WASM file from the builder stage to the current container
COPY --from=builder /project/plugin.wasm .