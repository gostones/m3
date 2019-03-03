#filebrowser
#https://docs.filebrowser.xyz/cli/filebrowser-config-set
define config_fb: method() = {
	# filebrowser config set \
	# 	--branding.name "File Browser" \
	# 	--port=50082 \
	# 	--auth.method=json
	filebrowser users update admin --password password
}

export FB_DATABASE = "${DHNT_BASE}/home/fb/filebrowser.db"

if ([ ! -e "${FB_DATABASE}" ]) {
	filebrowser config init \
		--branding.name "File Browser" \
		--port=50082 \
		--auth.method=json
	filebrowser users add admin password
}

config_fb

filebrowser --root ${DHNT_BASE}
printf "filebrowser exited"
    
