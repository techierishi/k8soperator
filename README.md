# k8soperator
A sample K8S Operator to demonstrate crud app deployment 

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).


# User Path

> Testing has been done on local kind cluster

## Deploy operator

```sh
# clone operator repo
git clone https://github.com/techierishi/k8soperator.git
cd k8soperator

# deploy to local cluster
helm install -f chart/values.yaml crud-helm ./chart/

# now create Crud kind
./bin/kustomize build config/samples | kubectl apply -f -
```

## Port forward

```sh
kubectl port-forward svc/mongocrud 8060:8060
```

## Get all resources:

```sh
kubectl get sc; kubectl get pv; kubectl get pvc; kubectl get pod; kubectl get service 
```

## Delete all created resources:

```sh
./bin/kustomize build config/samples | kubectl delete -f -; kubectl delete sc crud-storage-class; kubectl delete pv crud-pv
```

## Code explanation

`internal/controller/crud_controller.go` is the entrypoint of all reconcilers
`internal/controller/storage_class.go` is to create storage class with `AllowVolumeExpansion` capability
`internal/controller/pc.go` is to create persistent volume
`internal/controller/pvc.go` is to create persistent volume claim for database
`internal/controller/service.go` is to create node port service
`internal/controller/deployment.go` is to create deployment
`config/samples/schedule_v1_crud.yaml` has sample Crud kind yaml


---

# Developer Path

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=ghcr.io/techierishi/k8soperator:latest
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=ghcr.io/techierishi/k8soperator:latest
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
https://github.com/techierishi

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)


## Bundle commands

```sh
make bundle IMG=ghcr.io/techierishi/k8soperator:latest
make bundle-build BUNDLE_IMG=ghcr.io/techierishi/k8soperator-bundle:latest
make docker-push IMG=ghcr.io/techierishi/k8soperator-bundle:latest

```

## Generate helm charts

```sh
make helm
```

## Caveats

- Some cluster role has too much of permission make sure operator has enough role. This should be avoided in prod 
- Best coding practices may be missing at some places keeping the mind that this is a test app


## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

