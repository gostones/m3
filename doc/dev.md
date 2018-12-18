# Developer Guide

## Install required tools

*Docker CE/Kubernetes*

https://docs.docker.com/install/
https://kubernetes.io/docs/tasks/tools/install-kubectl/
https://docs.helm.sh/using_helm/#installing-helm

Start Docker CE and enable Kubernetes,
test with kubectl and helm with charts at https://github.com/helm/charts


*Git*

https://git-scm.com/book/en/v2/Getting-Started-Installing-Git



*IPFS*

https://github.com/ipfs/go-ipfs

Run 
`ipfs daemon`

Clone this repo and run
` bin/mirr`


*Chrome Proxy Plugin*

https://github.com/FelisCatus/SwitchyOmega/wiki/FAQ

Install SwitchyOmega or other proxy plugins of your choice
and point it to localhost:18080

This plugin is optional if you change your system default to localhost:18080


## Best practices 

In my humble opinion, programming language is not English so the latter's grammar rule does not apply to coding in the former.

For naming conventions:

prefer singular over plural
use word stem instead of its inflected variant


#### Principle
https://en.wikipedia.org/wiki/KISS_principle
https://en.wikipedia.org/wiki/Don%27t_repeat_yourself


#### Practice
https://peter.bourgon.org/go-best-practices-2016/
https://github.com/golang/go/wiki/CodeReviewComments
https://github.com/golang-standards/project-layout
