
```
kubectl apply -k deploy/config/manifests
```

```
cat << EOF | kubectl apply -f -
apiVersion: addon.open-cluster-management.io/v1alpha1
kind: ManagedClusterAddOn
metadata:
  name: clusternet
  namespace: cluster1
spec: {}
EOF
```