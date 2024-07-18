FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine as builder
WORKDIR /app
COPY . .

RUN set -ex; \
  dotnet publish -c Release; \
  mv /app/MDrop.Broker/bin/Release/net8.0/linux-x64/publish /app/;

FROM docker.io/library/alpine:latest as runner
COPY --from=builder /app/publish /app
WORKDIR /app

RUN set -ex; \
  apk update; \
  apk add gcc gcompat icu-dev icu-libs musl-dev tzdata; \
  rm appsettings.*.json;

CMD [ "/app/MDrop.Broker" ]
