FROM node:20-alpine AS builder

WORKDIR /workspace

RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/admin/package.json ./apps/admin/package.json

RUN pnpm install --frozen-lockfile

COPY apps/admin ./apps/admin

RUN pnpm -C apps/admin build

FROM node:20-alpine

WORKDIR /app

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

COPY --from=builder /workspace/apps/admin/.output /app/.output

EXPOSE 3000

CMD ["node", "/app/.output/server/index.mjs"]
