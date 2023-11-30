FROM golang:1.21.4-bullseye AS build

RUN apt-get update && apt-get install -y --no-install-recommends nkf
SHELL ["/bin/bash", "-c"]

WORKDIR /data
RUN set -o pipefail && \
  curl -Ls https://github.com/taku910/mecab/raw/master/mecab-ipadic/Noun.csv | \
  nkf -w | cut -f1,12 -d, > Noun.csv

ADD . /src
WORKDIR /src
RUN GOOS=linux go build -o wordchain main.go

FROM alpine
WORKDIR /app
COPY --from=build /src/wordchain /app
COPY --from=build /data/Noun.csv /app

ENTRYPOINT /app/wordchain
