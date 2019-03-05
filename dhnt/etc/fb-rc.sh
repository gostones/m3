#filebrowser
#https://docs.filebrowser.xyz/cli/filebrowser-config-set
define config_fb: method() = {
	# filebrowser config set \
	# 	--branding.name "M3 OS" \
	# 	--port=50082 \
	# 	--auth.method=json
	filebrowser users update dhnt --password password
}

export FB_DATABASE = "${DHNT_BASE}/home/fb/filebrowser.db"

if ([ ! -e "${FB_DATABASE}" ]) {
	filebrowser config init \
		--branding.name "M3 OS" \
		--port=50082 \
		--auth.method=json
	filebrowser users add dhnt password
}

config_fb

filebrowser --root ${DHNT_BASE}
printf "filebrowser exited"
    
