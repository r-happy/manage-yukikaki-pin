# syntax=docker/dockerfile:1

FROM node:20-bookworm AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

ENV NG_CLI_ANALYTICS="false"
RUN npm run build -- --configuration production

FROM node:20-bookworm-slim AS runner

WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

COPY --from=builder /app/dist ./dist

ENV NODE_ENV="production"
ENV PORT="4000"
ENV BACKEND_URL="http://backend:1323"

EXPOSE 4000

CMD ["node", "dist/frontend/server/server.mjs"]
