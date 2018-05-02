# kubectl plugin kill
Example kubectl plugin that force deletes a pod, removing any finalizers
that may be blocking deletion.

## Install the plugin
```console
$ curl -sLO https://github.com/carolynvs/kubectl-kill-plugin/releases/download/latest/kill.zip

$ unzip kill.zip -d ~/.kube/plugins/
Archive:  kill.zip
   creating: ~/.kube/plugins/kill/
  inflating: ~/.kube/plugins/kill/plugin.yaml
  inflating: ~/.kube/plugins/kill/kill
```

## Try the plugin

```console
$ make test
kubectl create -f test.yaml
pod "hello-world" created

kubectl get pod hello-world -o jsonpath='{.metadata.finalizers}'
[finalizer.kubernetes.io/hello-world]

kubectl plugin kill hello-world -v=0 --grace-period=0
removed finalizers from pod hello-world
killing default/hello-world with a grace period of 0s...
deleted pod hello-world

```

## Build the plugin
```console
$ make deploy
mkdir -p ~/.kube/plugins/kill
go build -o ~/.kube/plugins/kill/kill
cp plugin.yaml ~/.kube/plugins/kill/

```
