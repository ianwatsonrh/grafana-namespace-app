# Development
```
oc create sa cluster-reader
oc adm policy add-cluster-role-to-user cluster-reader -z cluster-reader
oc process -f gotemplate.json -p SOURCE_REPOSITOR_URL -p NAMESPACE< -p CLUSTER_READER_SERVICE__ACCOUNT | oc appy -f- 

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
oc create -f grafana-config.json
oc create -f grafana-dashboards.json
oc get secrets grafana-datasources -o json --export -n openshift-monitoring > grafana-datasources.json
oc create -f grafana-datasources.json
oc adm policy add-cluster-role-to-user system:auth-delegator -z grafana

```


# Invoke
oc create user test
oc adm add-role-to-user view test
curl http://admin-app-iw.apps.cacb.example.opentlc.com/search -X POST -H "Content-Type: application/json" --data '{"username":"test"}' -k
 
