# syntax=docker/dockerfile:1

FROM node:20-bookworm

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 4200

ENV NG_CLI_ANALYTICS="false"
ENV BACKEND_URL="https://api.mng-yukikaki.rhappy.dev"

CMD ["npm", "run", "start", "--", "--host", "0.0.0.0", "--port", "4200"]