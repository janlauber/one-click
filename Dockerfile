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
FROM node:lts-slim as ui-builder
WORKDIR /build
COPY ./frontend/package*.json ./
RUN rm -rf ./node_modules ./build
COPY ./frontend .
ARG APP_VERSION=0.0.1
RUN echo "VITE_APP_VERSION=$APP_VERSION" > .env
RUN npm install --legacy-peer-deps
RUN npm run build

# Stage 3: Runtime
FROM alpine as runtime
WORKDIR /app/one-click
COPY --from=backend-builder /build/one-click /app/one-click/one-click
COPY ./pocketbase/pb_migrations ./pb_migrations
COPY --from=ui-builder /build/build /app/one-click/pb_public
EXPOSE 8090
CMD ["/app/one-click/one-click", "serve", "--http", "0.0.0.0:8090"]
