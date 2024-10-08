## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.

> [!NOTE]
> Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Installation

> [!CAUTION]  
> Operator-Controller depends on [cert-manager](https://cert-manager.io/). Running the following command
> may affect an existing installation of cert-manager and cause cluster instability.

The latest version of Operator Controller can be installed with the following command:

```bash
curl -L -s https://github.com/operator-framework/operator-controller/releases/latest/download/install.sh | bash -s
```

### Additional setup on Macintosh computers
On Macintosh computers some additional setup is necessary to install and configure compatible tooling.

#### Install Homebrew and tools
Follow the instructions to [installing Homebrew](https://docs.brew.sh/Installation) and then execute the following to install tools:

```sh
brew install bash gnu-tar gsed
```

#### Configure your shell
Modify your login shell's `PATH` to prefer the new tools over those in the existing environment.  This example should work either with `zsh` (in $HOME/.zshrc) or `bash` (in $HOME/.bashrc):

```sh
for bindir in `find $(brew --prefix)/opt -type d -follow -name gnubin -print`
do
  export PATH=$bindir:$PATH
done
```

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/operator-controller:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/operator-controller:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
To undeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md) for more information.

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out

Install the CRDs and the operator-controller into a new [KIND cluster](https://kind.sigs.k8s.io/):
```sh
make run
```
This will build a local container image of the operator-controller, create a new KIND cluster and then deploy onto that cluster.
This will also deploy the catalogd and cert-manager dependencies.

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make help` for more information on all potential `make` targets.

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html).
