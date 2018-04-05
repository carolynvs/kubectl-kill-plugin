# Deploy the plugin
```console
mkdir -p ~/.kube/plugins/kill
go build -o ~/.kube/plugins/kill/kill
cp plugin.yaml ~/.kube/plugins/kill/
```

# Try the plugin

```console
$ kubectl plugin kill --help
Remove any finalizers on a pod, and delete it

Options:
      --grace-period='': Period of time in seconds given to the resource to terminate gracefully

Usage:
  kubectl plugin kill [flags] [options]

$ kubectl run hello-world --image=hello-world
deployment "hello-world" created

$ POD=$(kubectl get pods -o jsonpath='{.items[0].metadata.name}')

$ kubectl plugin kill $POD
removed finalizers from pod hello-world-ffbf4c44d-rs2tl
deleted pod hello-world-ffbf4c44d-rs2tl
```
