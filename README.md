# ENIConfig Controller

This repository will inplement auto annotating your Kubernete data plane nodes
with a desired `ENIConfig` name. This was originally implemented in the
`amazon-vpc-cni-k8s` project in this pull request - https://github.com/aws/amazon-vpc-cni-k8s/pull/165

## Running

```sh
# assumes you have a working kubeconfig, not required if operating in-cluster
$ go build -o eniconfig-controller .
$ ./eniconfig-controller -kubeconfig=$HOME/.kube/config -eniconfig-name=name-of-eni

## Deploying

If you would like to use this in-cluster you can modify the `args` in the
`configs/eniconfig-controller.yaml` to reference the name of the `ENIConfig`.


```yaml
args:
- --eniconfig-name=<name-of-eniconfig>
```

Then `apply` this config to the cluster.

```bash
kubectl apply -f eniconfig-controller.yaml
```

## Releasing

To release this project all you have to do is run `make release`.

## License

This library is licensed under the Apache 2.0 License.
