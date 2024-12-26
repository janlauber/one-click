# Stage 1: Build backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /build
COPY pocketbase/go.mod pocketbase/go.sum pocketbase/main.go ./
COPY pocketbase/pkg ./pkg
RUN apk --no-cache add upx make git gcc libtool musl-dev ca-certificates dumb-init \
  && go mod tidy \
  && CGO_ENABLED=0 go build \
  && upx one-click

# Stage 2: Build frontend
FROM node:lts-slim AS ui-builder
WORKDIR /build
COPY ./frontend/package*.json ./
RUN rm -rf ./node_modules ./build
COPY ./frontend .
ARG APP_VERSION=0.0.1
RUN echo "VITE_APP_VERSION=$APP_VERSION" > .env
RUN npm install --legacy-peer-deps
RUN npm run build

# Stage 3: Runtime
FROM alpine AS runtime
WORKDIR /app/one-click

# Install ca-certificates package and create directory for custom certificates
RUN apk --no-cache add ca-certificates \
    && mkdir -p /usr/local/share/ca-certificates

# Copy application files
COPY --from=backend-builder /build/one-click /app/one-click/one-click
COPY ./pocketbase/pb_migrations ./pb_migrations
COPY --from=ui-builder /build/build /app/one-click/pb_public

# Create entrypoint script to handle certificates
COPY <<EOF /entrypoint.sh
#!/bin/sh
# If CUSTOM_CA_CERT environment variable is set, add the certificate
if [ -n "\${CUSTOM_CA_CERT}" ]; then
    echo "\${CUSTOM_CA_CERT}" > /usr/local/share/ca-certificates/custom-ca.crt
    update-ca-certificates
fi

# Execute the main application
exec /app/one-click/one-click serve --http "0.0.0.0:8090"
EOF

RUN chmod +x /entrypoint.sh

EXPOSE 8090
ENTRYPOINT ["/entrypoint.sh"]
