FROM golang:1.17-alpine

WORKDIR /go/src/
COPY . . 
RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk add --no-cache bash git openssh make cmake
EXPOSE 4318/tcp
EXPOSE 4318/udp
ENV LOGICMONITOR_ACCOUNT=$LOGICMONITOR_ACCOUNT
ENV LOGICMONITOR_BEARER_TOKEN=$LOGICMONITOR_BEARER_TOKEN
RUN go build -o otel-collector-action
CMD ["/go/src/otel-collector-action"]
