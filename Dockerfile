ARG NODE_VERSION=24
ARG PNPM_VERSION=10.32.1
ARG GO2RTC_IMAGE=alexxit/go2rtc:latest

FROM node:${NODE_VERSION}-alpine AS base
RUN apk add --no-cache libc6-compat
ARG PNPM_VERSION
RUN corepack enable && corepack prepare pnpm@${PNPM_VERSION} --activate
WORKDIR /app

FROM base AS build-deps
RUN apk add --no-cache python3 make g++
COPY pnpm-workspace.yaml package.json pnpm-lock.yaml ./
COPY server/package.json ./server/
COPY web/package.json ./web/
RUN pnpm install --frozen-lockfile

FROM build-deps AS web-builder
COPY web ./web
RUN pnpm -F web build

FROM build-deps AS server-builder
COPY server ./server
RUN pnpm -F server build

FROM base AS prod-deps
RUN apk add --no-cache python3 make g++
COPY pnpm-workspace.yaml package.json pnpm-lock.yaml ./
COPY server/package.json ./server/
COPY web/package.json ./web/
RUN pnpm install --frozen-lockfile --prod --filter server...

FROM ${GO2RTC_IMAGE} AS go2rtc

FROM node:${NODE_VERSION}-alpine AS runtime
RUN apk add --no-cache tini libc6-compat bash
WORKDIR /app

COPY --from=go2rtc /usr/local/bin/go2rtc /usr/local/bin/go2rtc

COPY --from=server-builder /app/server/dist ./server/dist
COPY --from=server-builder /app/server/package.json ./server/package.json
COPY --from=prod-deps /app/server/node_modules ./server/node_modules
COPY --from=prod-deps /app/node_modules ./node_modules

COPY --from=web-builder /app/web/build ./web

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh && mkdir -p /config /data

ENV NODE_ENV=production
ENV PORT=3000
ENV WEB_STATIC_PATH=/app/web
ENV GO2RTC_WS_URL=ws://localhost:1984

EXPOSE 3000 1984 8555

WORKDIR /app/server
ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
