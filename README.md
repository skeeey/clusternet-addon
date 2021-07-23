1. Deploy clusternet-hub on the OCM hub cluster

2. Deploy the clusternet-addon hub controller on the OCM hub cluster
```
kubectl apply -k deploy/config/manifests
```

3. Create a clusternet ManagedClusterAddOn in the OCM managed cluster namespace (e.g. `cluster1`) on the OCM hub cluster
```
cat << EOF | kubectl apply -f -
apiVersion: addon.open-cluster-management.io/v1alpha1
kind: ManagedClusterAddOn
metadata:
  name: clusternet
  namespace: cluster1
spec:
  installNamespace: clusternetaddon
EOF
```

4. Create a clusternet ManagedCluster in the OCM managed cluster namespace (e.g. `cluster1`) on the OCM hub cluster
```
cat << EOF | kubectl apply -f -
apiVersion: clusters.clusternet.io/v1beta1
kind: ManagedCluster
metadata:
  labels:
    clusters.clusternet.io/cluster-id: cluster1
  name: clusternet-cluster1
  namespace: cluster1
spec:
  clusterId: 00000000-0000-0000-0000-000000000000
  syncMode: Dual
EOF
```