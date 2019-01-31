#!/usr/bin/env bash

# https://gtt.dst.ibm.com/tools/vacationplanner/calendar.php?team=17782&year=2014&month=5&display=month
# https://docs.traefik.io/user-guide/kv-config/#upload-the-configuration-in-the-key-value-store

# docker run \
#   -v /var/run/docker.sock:/var/run/docker.sock \
#   -v $DHNT_BASE/etc/traefik:/etc/traefik \
#   -p 38080:80 \
#   -p 38081:81 \
#   -p 38443:443 \
#   -l traefik.frontend.rule=Host:traefik.home \
#   -l traefik.port=80 \
#   --network web \
#   --name traefik \
#   traefik:1.7.2-alpine

docker inspect web; if [ $? -ne 0 ]; then
    echo "creating network web ..."
    docker network create web
fi

docker run \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $PWD/etc/traefik:/etc/traefik \
  -p 38080:80 \
  -p 38081:8080 \
  -p 38443:443 \
  --network web \
  --name traefik \
  traefik:1.7.2-alpine