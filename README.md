# sample-crd

This is an example of kubernetes CRD.It will create a deployment which name is defined in SampleCRD.spec.deploymentName

```
kubectl create -f examples/crd.yaml

go run *.go --kubeconfig=$HOME/.kube/config

kubectl create -f examples/samplecrd.yaml

kubectl get samplecrds,deployment,pod

```