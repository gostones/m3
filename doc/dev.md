# Developer Guide

## Set up


You may run setup.sh to have all the required projects set up automatically as follows:

```
wget -O - https://raw.githubusercontent.com/dhnt/m3/master/setup.sh | bash

#
<!-- source ~/dhnt/m3/env.sh -->

sudo m3d install --base ~/dhnt
sudo m3d start

systray

sudo m3d stop
sudo m3d uninstall

```

If everything goes well, you should now have M3 up and running.

However if it failed, please follow the following manual steps.

Assuming M3 source code and third party projects will be checked out under `dhnt`.
For third party Golang projects that still require the GOPATH to build will go into dhnt/go.



```
#mkdir -p ~/dhnt
#mkdir -p ~/dhnt/go
mkdir -p ~/dhnt/go/bin

export DHNT_BASE=~/dhnt
export GOPATH=~/dhnt/go
export PATH=./:$GOPATH/bin:$DHNT_BASE/bin:$PATH
```

## Install required tools

*Docker CE/Kubernetes*

https://docs.docker.com/install/
https://kubernetes.io/docs/tasks/tools/install-kubectl/
https://docs.helm.sh/using_helm/#installing-helm

Start Docker CE and enable Kubernetes,
test with kubectl and helm with charts at https://github.com/helm/charts


*Git*

Client

https://git-scm.com/book/en/v2/Getting-Started-Installing-Git

Git server - Gogs

https://github.com/gogs/gogs

Install and run from [source](https://gogs.io/docs/installation/install_from_source.html)

```
export GOPATH=~/dhnt/go

mkdir -p $GOPATH/src/github.com/gogs
cd $GOPATH/src/github.com/gogs
git clone https://github.com/gogs/gogs.git
cd gogs
go build -tags "sqlite pam cert"

./gogs web
```

Visit: http://0.0.0.0:3000


*IPFS*

https://github.com/ipfs/go-ipfs

build and Run from [source](https://github.com/ipfs/go-ipfs#development)

```
export GOPATH=~/dhnt/go
export PATH=$GOPATH/bin:$PATH

mkdir $GOPATH/src/github.com/ipfs
cd $GOPATH/src/github.com/ipfs
git clone https://github.com/ipfs/go-ipfs.git
cd go-ipfs
make install

#optional - change default ports
ipfs config Addresses
ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/9001
ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001
#

ipfs daemon
```

*M3*

Clone this repo and run from [source](https://github.com/dhnt/m3.git)

```
cd ~/dhnt
git clone https://github.com/dhnt/m3.git
cd m3
./build.sh

```

run `pmd --port 18082`

ps command should list something similar to the following:

```
$ps

13751 ttys001    0:00.06 pmd -port 18082
13752 ttys001    0:00.07 /Users/liqiang/dhnt/go/bin/mirr --port 18080
13753 ttys001    0:00.03 /Users/liqiang/dhnt/go/bin/gotty --port 50022 --permit-wr
13754 ttys001    0:00.44 /Users/liqiang/dhnt/go/bin/gogs web --port 3000
13755 ttys001    0:00.19 /Users/liqiang/dhnt/go/bin/ipfs daemon
```

*Chrome Proxy Plugin*

https://github.com/FelisCatus/SwitchyOmega/wiki/FAQ

Install SwitchyOmega or other proxy plugins of your choice
and point it to localhost:18080

This plugin is optional if you change your system default to localhost:18080

## References

[Base32 Encoding](http://www.crockford.com/wrmg/base32.html)

## Best practices 

In my humble opinion, programming language is not English so the latter's grammar rule does not apply to coding in the former.

For naming conventions:

prefer singular over plural.

use word stem instead of its inflected variant.


#### Principle
[KISS](https://en.wikipedia.org/wiki/KISS_principle)

[DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)


#### Practice
https://peter.bourgon.org/go-best-practices-2016/

https://github.com/golang/go/wiki/CodeReviewComments

https://github.com/golang-standards/project-layout
