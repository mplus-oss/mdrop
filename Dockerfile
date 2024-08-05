# Build tunnel
FROM docker.io/library/golang:latest AS build
WORKDIR /build
COPY . .

RUN make build-tunnel; \
  chmod +x mdrop-tunnel;

# Create Private Environment
FROM docker.io/library/alpine:latest
LABEL org.opencontainers.image.authors="Ikramullah <ikramullah@mplus.software>,Syahrial Agni Prasetya <syahrial@mplus.software>"
WORKDIR /

RUN set -ex; \
  apk update; \
  apk add openssh bash lsof openssl;

COPY ./tunnel.conf /etc/ssh/sshd_config.d/
COPY ./entrypoint.sh .
COPY ./mdrop-tunnel.sh .
COPY --from=build /build/mdrop-tunnel /usr/bin/mdrop-tunnel

RUN adduser tunnel -HD; \
  passwd tunnel -d; \
  ssh-keygen -A;

EXPOSE 22

ENTRYPOINT [ "/entrypoint.sh" ]
