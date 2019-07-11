# Development
```
oc create sa cluster-reader
oc adm policy add-cluster-role-to-user cluster-reader -z cluster-reader
oc process -f openshift/template.json | oc apply -f- 

#datasource.js - var interpolated = {"username": this.contextSrv.user.login}; then add contextSrv into constructor

oc import-image devtools/go-toolset-rhel7 --from=registry.access.redhat.com/devtools/go-toolset-rhel7 --confirm


Either SSH key or Token
#ssh
oc create secret generic scmsecret --from-file=ssh-privatekey=$HOME/.ssh/id_rsa --dry-run -o json | oc apply -f -
oc annotate secret scmsecret 'build.openshift.io/source-secret-match-uri-1=ssh://github.com/ianwatsonrh/*'


#token
oc create secret generic scmsecret --from-literal=password=Value --type=kubernetes.io/basic-auth
oc annotate secret scmsecret 'build.openshift.io/source-secret-match-uri-1=https://github.com/ianwatsonrh/*'


#grafana
oc create serviceaccount grafana
oc create secret generic grafana-proxy --from-literal=session_secret=$(openssl rand -base64 13)
oc create secret generic grafana-config --from-file=openshift/grafana.ini
oc create -f openshift/grafana-dashboards.json
oc create secret generic grafana-datasources --from-file=openshift/grafana-datasources.yaml
oc adm policy add-cluster-role-to-user system:auth-delegator -z grafana
oc create configmap grafana-dashboard-rbac-example --from-file=openshift/rbac-example.json

POD=$(oc get pods | grep grafana | awk '{print $1}')
oc rsync ./grafana-extension/simple-json-datasource-master $POD:/var/lib/grafana/plugins/
oc delete pod $POD
```


# Invoke example
htpasswd -b /etc/origin/master/htpasswd test test
oc create user test
oc adm policy add-role-to-user view test
curl http://admin-app-iw.apps.cacb.example.opentlc.com/search -X POST -H "Content-Type: application/json" --data '{"username":"test"}' -k
 
