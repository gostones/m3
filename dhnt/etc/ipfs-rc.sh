#
define config_ipfs: method() = {
	define myid = `(ipfs id "--format=<id>") #`vscode formmatting hack`
	echo "configuring $myid ..."
	#ipfs config Addresses
	#optional - change default ports
	ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/58080 #8080
	ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001
	ipfs config --json Experimental.Libp2pStreamMounting true
	ipfs config --json Experimental.FilestoreEnabled true
	ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["http://ipfs.home", "http://127.0.0.1:5001", "https://webui.ipfs.io"]'
	ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
}

export IPFS_PATH = "${DHNT_BASE}/home/ipfs"

if ([ ! -e "${IPFS_PATH}/config" ]) {
	ipfs init
}
config_ipfs

ipfs daemon
printf "ipfs exited"
    
