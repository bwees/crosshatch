# Production image, assembled by GoReleaser (dockers_v2).
# The server binary is built by GoReleaser per-platform and the web bundle is
# built in CI; both are provided via the build context, so this only assembles
# the runtime layer.
ARG GO2RTC_IMAGE=alexxit/go2rtc:latest

FROM ${GO2RTC_IMAGE} AS go2rtc

FROM alpine:3.20 AS runtime
ARG TARGETPLATFORM
RUN apk add --no-cache tini bash ca-certificates libc6-compat
WORKDIR /app

COPY --from=go2rtc /usr/local/bin/go2rtc /usr/local/bin/go2rtc
COPY ${TARGETPLATFORM}/crosshatch /app/server/crosshatch
COPY web/build ./web
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /app/server/crosshatch && mkdir -p /config /data

ENV PORT=3000
ENV WEB_STATIC_PATH=/app/web
ENV GO2RTC_WS_URL=ws://localhost:1984
ENV DATA_DIR=/data

EXPOSE 3000 1984 8555

WORKDIR /app/server
ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
