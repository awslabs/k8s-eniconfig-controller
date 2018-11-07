# ENIConfig Controller

This repository will inplement auto annotating your Kubernete data plane nodes
with a desired `ENIConfig` name. This was originally implemented in the
`amazon-vpc-cni-k8s` project in this pull request - https://github.com/aws/amazon-vpc-cni-k8s/pull/165

## Deploying

If you use `helm` to deploy applications you can use the
`charts/eniconfig-controller` chart to deploy the controller.

```bash
helm install --set eniConfigName=<name-of-eniconfig-cr> ./charts/eniconfig-controller
```

If you do not use `helm` you can download the config file using the below steps
`wget` or `curl -o` the `configs/eniconfig-controller.yaml` then modify the
`args` to represent the proper `ENIConfig` name. E.G.

```yaml
args:
- --eniconfig-name=<name-of-eniconfig>
```

Then `apply` this config to the cluster.

```bash
kubectl apply -f eniconfig-controller.yaml
```

## Running in Dev

```sh
# assumes you have a working kubeconfig, not required if operating in-cluster
$ go build -o eniconfig-controller .
$ ./eniconfig-controller -kubeconfig=$HOME/.kube/config -eniconfig-name=name-of-eni

## Releasing

To release this project all you have to do is run `make release`.
## License

This library is licensed under the Apache 2.0 License.
