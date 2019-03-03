#
export GOGS_WORK_DIR = "${DHNT_BASE}/home/gogs"
export USER = git

cd ${DHNT_BASE}
gogs web --port 3000
##