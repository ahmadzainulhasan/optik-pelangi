FROM oven/bun:canary-alpine AS builder

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run --bun build

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache libstdc++ && \
	addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/build/server /app/server

USER appuser

CMD ["/app/server"]
