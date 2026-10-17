FROM oven/bun:canary-alpine AS builder

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run --bun build

FROM oven/bun:canary-alpine

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/build/server /app/server

USER appuser

CMD ["/app/server"]
