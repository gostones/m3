Mirr
-----

An IPFS DApp that mirrors the old web, a.k.a,  the world wide web.
it has a built-in forward proxy and load balancer to distribute the load to multiple peers.


### Usage

```
$ ./mirr -port 18080
```

### Credits

https://github.com/voldyman/GoLoadBalance

https://github.com/kintoandar/fwd

https://github.com/elazarl/goproxy

<!-- https://github.com/FelisCatus/SwitchyOmega -->
<!-- https://github.com/PuerkitoBio/gocrawl -->
<!-- https://github.com/gocolly/colly -->

### License

Mirr is released under MIT license

Author: Qiang Li <liqiang@gmail.com>

<!--
https://docs.ipfs.io/reference/api/http/

curl "http://localhost:5001/api/v0/swarm/addrs/local?id=<value>"
curl http://127.0.0.1:5001/api/v0/swarm/peers

curl "http://localhost:5001/api/v0/p2p/stream/dial?arg=<Peer>&arg=<Protocol>&arg=<BindAddress>"

-->


<!-- 
127.0.0.1
::1
localhost
-->

<!-- 

# UI
http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/overview?namespace=_all


#traefik MIT 19,017 Go
https://github.com/containous/traefik

helm install stable/traefik --name traefik --namespace kube-system \
	--set ssl.insecureSkipVerify=true \
	--set dashboard.enabled=true \
	--set dashboard.domain=localhost

helm upgrade traefik stable/traefik \
	--set ssl.insecureSkipVerify=true \
	--set dashboard.enabled=true \
	--set dashboard.domain=localhost


#kubectl apply -f ./traefik-ui.yaml
#https://github.com/containous/traefik/blob/master/examples/k8s/ui.yaml
https://docs.traefik.io/user-guide/kubernetes/
https://github.com/helm/charts/tree/master/stable/traefik

###cicd:

#gogs MIT 28,188 Go
https://github.com/gogs/gogs

helm install --namespace cicd --name gogs incubator/gogs \
	--set serviceType=ClusterIP


#gitlab MIT  21,405 Ruby
https://github.com/gitlabhq/gitlabhq

helm install --namespace cicd --name gitlab stable/gitlab-ce \
	--set serviceType=ClusterIP \
	--set externalUrl=http://1220490149ec3a5ccf6ac3d8db2ec7c42e8486b7e95c0a324a0eaf22ae50d2fc1011/


 http://1220490149ec3a5ccf6ac3d8db2ec7c42e8486b7e95c0a324a0eaf22ae50d2fc1011/

Username: root
Password: <whatever value you entered

#jenkins MIT 11,671 Java
https://github.com/jenkinsci/jenkins

helm install --namespace cicd --name jenkins stable/jenkins

export SERVICE_IP=$(kubectl get svc --namespace cicd jenkins --template "{{ range (index .status.loadBalancer.ingress 0) }}{{ . }}{{ end }}")
echo http://$SERVICE_IP:8080/login
admin


#sonarqube LGPL3 3,151 Java
https://github.com/SonarSource/sonarqube 

helm install --namespace cicd --name sonarqube stable/sonarqube

export SERVICE_IP=$(kubectl get svc --namespace cicd sonarqube-sonarqube -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
http://$SERVICE_IP:9000

#gocd Apache 4,491 Java

https://github.com/theia-ide/theia
https://github.com/b3log/wide



###cwe
https://en.wikipedia.org/wiki/Collaborative_working_environment


#mattermost MIT 13,623 Go 
helm install --namespace cwe --name mattermost stable/mattermost-team-edition \
  --set mysql.mysqlUser=admin \
  --set mysql.mysqlPassword=password 

kubectl port-forward --namespace cwe $(kubectl get pods --namespace cwe -l "app=mattermost-mattermost-team-edition,release=mattermost" -o jsonpath='{ .items[0].metadata.name }') 8080:8065
helm install --set host=mattermost.yourdomain.com --set ingress.enabled=true stable/mattermost-team-edition


#wordpress GNU 11,707 PHP
https://github.com/WordPress/WordPress

helm install --namespace cwe --name wordpress stable/wordpress \
	--set service.type=ClusterIP \
	--set wordpressUsername=admin,wordpressPassword=password \
	--set mariadb.mariadbRootPassword=secretpassword

echo Username: admin
echo Password: $(kubectl get secret --namespace cwe wordpress-wordpress -o jsonpath="{.data.wordpress-password}" | base64 --decode)


#mediawiki GNU  1,364 PHP
https://github.com/wikimedia/mediawiki

helm install  --namespace cwe  --name mediawiki stable/mediawiki \
	--set service.type=ClusterIP 

kubectl port-forward --namespace cwe svc/mediawiki-mediawiki 18082:80
echo Username: user
echo Password: $(kubectl get secret --namespace cwe mediawiki-mediawiki -o jsonpath="{.data.mediawiki-password}" | base64 --decode)


#dokuwiki GPL 2,317 PHP
https://github.com/splitbrain/dokuwiki

helm install --namespace cwe --name dokuwiki stable/dokuwiki \
	--set service.type=ClusterIP 

kubectl port-forward --namespace cwe svc/dokuwiki-dokuwiki 18081:80
echo Username: user / Q6le5NSVwT
echo Password: $(kubectl get secret --namespace cwe dokuwiki-dokuwiki -o jsonpath="{.data.dokuwiki-password}" | base64 --decode)


-->