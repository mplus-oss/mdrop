FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine as build
LABEL org.opencontainers.image.authors="Ikramullah <ikramullah@mplus.software>"
WORKDIR /app
COPY . .

RUN set -ex; \
  dotnet restore; \
  dotnet publish --no-restore -c Release -o ./Application

FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine-composite
LABEL org.opencontainers.image.authors="Ikramullah <ikramullah@mplus.software>"
WORKDIR /app

COPY --from=build /app/Application /app
COPY --from=build /app/*.sh /app/

RUN set -ex; \
  apk update; \
  apk add openssl;

ENTRYPOINT [ "/app/entrypoint.sh" ]
