FROM golang:latest

RUN apt update && apt install -y chromium

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main .

CMD [ "/app/main" ]