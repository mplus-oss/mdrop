FROM docker.io/library/alpine:latest

RUN set -ex; \
  apk update; \
  apk add openssh bash lsof;

COPY ./start.sh /
COPY ./tunnel.conf /etc/ssh/sshd_config.d/
COPY ./entrypoint.sh /

RUN adduser tunnel -HD; \
  passwd tunnel -d; \
  ssh-keygen -A;

EXPOSE 22

ENTRYPOINT /entrypoint.sh
