FROM golang:1.20 AS build

ARG USER="platform-go.ro-bot"
ARG ACCESS_TOKEN="Ms1xVaFWvjgBsB4dLo-j"

WORKDIR /go/src/app

RUN go env -w GOPRIVATE=gitlab.tubecorporate.com

RUN git config \
    --global \
    url."https://${USER}}:${ACCESS_TOKEN}@gitlab.tubecorporate.com".insteadOf \
    "https://gitlab.tubecorporate.com"

# Cache dependencies
ADD go.mod .
ADD go.sum .


RUN go mod download

ADD . .

RUN go install github.com/magefile/mage && \
    mage build

FROM alpine:3.17 AS app

RUN apk --no-cache add \
        ca-certificates \
        curl \
        unzip \
        make \
        wget \
        htop \
        net-tools \
        curl \
        tzdata


WORKDIR /usr/bin/reminder

COPY --from=build /go/src/app/build /usr/bin/reminder/
COPY --from=build /go/src/app/config /usr/bin/reminder/

ENTRYPOINT /usr/bin/reminder/reminder
