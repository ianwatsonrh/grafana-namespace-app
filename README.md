# Development
oc create sa cluster-reader
oc adm policu add-cluster-role-to-user cluster-reader -z cluster-reader
oc process -f gotemplate.json -p SOURCE_REPOSITOR_URL -p NAMESPACE< -p CLUSTER_READER_SERVICE__ACCOUNT | oc appy -f- 

datasource.js - var interpolated = {"username": this.contextSrv.user.login}; then add contextSrv into constructor


oc create secret generic scmsecret --from-file=ssh-privatekey=$HOME/.ssh/id_rsa --dry-run -o json | oc apply -f -
oc annotate secret scmsecret 'build.openshift.io/source-secret-match-uri-1=ssh://github.com/ianwatsonrh/*'


 
