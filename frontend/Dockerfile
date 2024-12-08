FROM node:22.11-alpine AS build

WORKDIR /app
COPY package.json package-lock.json* ./

RUN npm install

COPY . .

RUN npm run build

FROM node:22.11-alpine

WORKDIR /app

RUN npm install -g serve

COPY --from=build /app/dist ./dist

ENV NODE_ENV=production

EXPOSE 5000

CMD ["serve", "-s", "dist", "-l", "5000"]
