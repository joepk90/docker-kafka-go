# =========================
# STAGE 1: build go app
# =========================
FROM golang:1.23 AS build

# ARG GITHUB_TOKEN
# RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

ADD . /app
WORKDIR /app
RUN make dep
RUN make build

# =========================
# STAGE 3: runtime
# =========================
FROM alpine:3.17
RUN addgroup -g 1000 -S appgroup && adduser -u 1000 -S appuser -G appgroup

RUN mkdir /app
COPY --from=build /app/docker-kafka-go /app/docker-kafka-go
USER appuser

ENTRYPOINT exec /app/docker-kafka-go
