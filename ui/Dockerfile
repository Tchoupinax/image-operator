FROM node:22-alpine

LABEL buildAt=2024-11-02

ARG JSON

WORKDIR /app
USER 1001

COPY .output /app/.output
COPY package.json /app/package.json

CMD HOST=0.0.0.0 PORT=3000 node /app/.output/server/index.mjs
