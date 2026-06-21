ARG NODE_VERSION=24
ARG PNPM_VERSION=10.32.1
ARG GO_VERSION=1.26
ARG GO2RTC_IMAGE=alexxit/go2rtc:latest

# ---- Web build (SvelteKit static) ----
FROM node:${NODE_VERSION}-alpine AS web-base
RUN apk add --no-cache libc6-compat
ARG PNPM_VERSION
RUN corepack enable && corepack prepare pnpm@${PNPM_VERSION} --activate
WORKDIR /app

FROM web-base AS web-deps
RUN apk add --no-cache python3 make g++
WORKDIR /app/web
COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM web-deps AS web-builder
COPY web ./
RUN pnpm build

# ---- Go server build ----
FROM golang:${GO_VERSION}-alpine AS server-builder
# build-base provides the C toolchain for the CGO-based sqlite driver.
RUN apk add --no-cache build-base
WORKDIR /src
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=1 go build -trimpath -o /out/crosshatch .

# ---- go2rtc binary ----
FROM ${GO2RTC_IMAGE} AS go2rtc

# ---- Runtime ----
FROM alpine:3.20 AS runtime
RUN apk add --no-cache tini bash ca-certificates libc6-compat
WORKDIR /app

COPY --from=go2rtc /usr/local/bin/go2rtc /usr/local/bin/go2rtc
COPY --from=server-builder /out/crosshatch /app/server/crosshatch
COPY --from=web-builder /app/web/build ./web

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh && mkdir -p /config /data

ENV PORT=3000
ENV WEB_STATIC_PATH=/app/web
ENV GO2RTC_WS_URL=ws://localhost:1984
ENV DATABASE_URL=/data/crosshatch.db

EXPOSE 3000 1984 8555

WORKDIR /app/server
ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
