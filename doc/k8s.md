# Home Node Setup

## K8S - Kubernetes

Open Docker > Preferences > Kubernetes

Check "Enable Kerbernetes"

## Helm Charts

*Official*

https://github.com/helm/charts

*Kompose*

http://kompose.io/


*create your own*

https://docs.helm.sh/developing_charts/


## CICD - continuous integration and deployment or delivery

*Gogs*

https://github.com/gogs/gogs

```
helm install --namespace cicd --name gogs incubator/gogs \
	--set serviceType=ClusterIP
```

*Jenkins*

https://github.com/jenkinsci/jenkins

```
helm install --namespace cicd --name jenkins stable/jenkins
```

*Sonarqube*

https://github.com/SonarSource/sonarqube 


```
helm install --namespace cicd --name sonarqube stable/sonarqube
```

## CWE - collaborative working environment

*Mattermost*

```
helm install --namespace cwe --name mattermost stable/mattermost-team-edition \
  --set mysql.mysqlUser=admin \
  --set mysql.mysqlPassword=password \
  --set config.SiteUrl=http://chat.home/ \
  --set ingress.enabled=false
  ```

  *Wordpress*

https://github.com/WordPress/WordPress

```
helm install --namespace cwe --name wordpress stable/wordpress \
	--set service.type=ClusterIP \
	--set wordpressUsername=admin,wordpressPassword=password \
	--set mariadb.mariadbRootPassword=secretpassword
```

*Dokuwiki*

https://github.com/splitbrain/dokuwiki

```
helm install --namespace cwe --name dokuwiki stable/dokuwiki \
	--set service.type=ClusterIP 
```



## Ingress Controller - reverse proxy

We use [Traefik](https://docs.traefik.io/) [Helm chart](https://github.com/helm/charts/tree/master/stable/traefik)

```
helm install stable/traefik --name traefik --namespace kube-system \
	--set ssl.insecureSkipVerify=true \
	--set cpuLimit=500m \
	--set memoryLimit=1Gi \
	--set dashboard.enabled=true \
	--set dashboard.domain=localhost
```

To allow incoming traffic to any applications you have deployed, enable \.home and \.\[hexid\] domains

To find your own peer id:
`ipfs id` and convert it to hex enocoded peer address (hexid) `bin/hexid --id your-peer-id`

Example (after you have successfully installed mattermost):

```
# kubectl -n cwe apply -f mattermost-ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: mattermost-ingress
  annotations:
    kubernetes.io/ingress.class: traefik
spec:
  rules:
  - host: chat.home
    http:
      paths:
      - path: /
        backend:
          serviceName: mattermost-team-edition
          servicePort: 8065
  - host: chat.<hexid>
    http:
      paths:
      - path: /
        backend:
          serviceName: mattermost-team-edition
          servicePort: 8065
```

Replace hexid in the above with your own, save as `mattermost-ingress.yaml` and run `kubectl -n cwe apply -f mattermost-ingress.yaml`

See more examples in internal/k8s and you can learn more here: https://kubernetes.io/docs/concepts/services-networking/ingress/


Verify with curl:

```
kubectl get ingress -n cwe
```

If you see your entries, you can verify:

```
curl -kv http://localhost:80 -H "Host: chat.home"
```

If you see response, verify via M3 proxy.

```
curl -kv -x http://localhost:18080 http://chat.home 

```

You should now be able to access mattermost from browser by entering http://chat.home in the address bar after you have installed a browser proxy plugin pointing to http://localhost:18080
