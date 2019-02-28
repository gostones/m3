#
cd ${DHNT_BASE}

traefik -c "${DHNT_BASE}/etc/traefik/traefik.toml" "--file.directory=${DHNT_BASE}/etc/traefik/config/"
##