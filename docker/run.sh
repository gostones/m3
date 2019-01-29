
docker run \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $DHNT_BASE/etc/traefik:/etc/traefik \
  -p 28080:28080 \
  -p 28081:28081 \
  -p 28443:28443 \
  -l traefik.frontend.rule=Host:traefik.home \
  -l traefik.port=28080 \
  --network web \
  --name traefik \
  traefik:1.7.2-alpine