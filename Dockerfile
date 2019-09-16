FROM golang:latest as builder

RUN apt update && apt install -y tzdata
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download 
COPY . .
RUN make build

FROM alpine:latest

ENV DOCKER_API_VERSION="1.40"
ENV TZ=Europe/Moscow
ENV REGISTRY_IP=""
ENV REGISTRY_PORT="5000"
ENV CRONTAB="0 0 0 * * *"
ENV LOG_LEVEL=ERROR
ENV APP_PREFIX=""
ENV PERIOD=60
ENV IMAGE_AMOUNT=5
ENV AUTOUPDATE=1
ENV REGISTRY_PATH="/var/lib/registry"

COPY --from=builder /app/bin/drwatcher /bin/drwatcher

CMD ["/bin/drwatcher"]