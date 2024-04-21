FROM node:latest

WORKDIR /chat

COPY . .

RUN yarn install

CMD npx prisma migrate dev --name init \
    && yarn start