# Development
oc create sa cluster-reader
oc adm policu add-cluster-role-to-user cluster-reader -z cluster-reader
oc process -f gotemplate.json -p SOURCE_REPOSITOR_URL -p NAMESPACE< -p CLUSTER_READER_SERVICE__ACCOUNT | oc appy -f- 
