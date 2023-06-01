FROM golang:latest as builder1
RUN mkdir /app
COPY . .
RUN go build -o app main.go

FROM ubuntu as builder2
RUN apt update && \
    apt install wget unzip && \
    mkdir /v2ray
RUN wget https://github.com/v2fly/v2ray-core/releases/download/v4.28.2/v2ray-linux-64.zip && \
    unzip v2ray-linux-64.zip && \
    rm -rf v2ray-linux-64.zip

FROM ubuntu
RUN mkdir /v2ray
WORKDIR /v2ray
COPY --from=builder2 /v2ray/ .
COPY --from=builder1 /app/app .

CMD ["./start.sh"]