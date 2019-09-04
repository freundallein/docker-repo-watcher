FROM golang:latest as builder

# RUN useradd -u 10001 scratchuser
RUN apt update && apt install -y tzdata
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download 

COPY . .
RUN make build

FROM scratch

ENV DOCKER_API_VERSION="1.40"
ENV TZ=Europe/Moscow
ENV REGISTRY_IP=""
ENV REGISTRY_PORT="5000"
ENV CRONTAB="0 0 * * *"
ENV LOG_LEVEL=ERROR
ENV APP_PREFIX=""
ENV PERIOD=60
ENV IMAGE_AMOUNT=5
ENV AUTOUPDATE=1


COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/bin/drwatcher /bin/drwatcher

# USER scratchuser
CMD ["/bin/drwatcher"]