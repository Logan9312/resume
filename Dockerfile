# Build stage
FROM rust:1.83-bookworm AS builder

WORKDIR /app

# Copy manifests first for dependency caching
COPY Cargo.toml Cargo.lock* ./

# Create dummy src to build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy actual source and rebuild
COPY src ./src
RUN touch src/main.rs && cargo build --release

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Download PDFium binary from bblanchon/pdfium-binaries
RUN curl -L -o pdfium.tgz https://github.com/ArtifexSoftware/pdfium-lib/releases/download/1.1/pdfium-linux-x64-release.tgz \
    && tar -xzf pdfium.tgz \
    && mv lib/libpdfium.so . \
    && rm -rf pdfium.tgz lib include

COPY --from=builder /app/target/release/resume .
COPY resume.pdf .

ENV LD_LIBRARY_PATH=/app

EXPOSE 3000

CMD ["./resume"]
