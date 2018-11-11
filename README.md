# ENIConfig Controller

This repository will inplement auto annotating your Kubernete data plane nodes
with a desired `ENIConfig` name. This was originally implemented in the
`amazon-vpc-cni-k8s` project in this pull request - https://github.com/aws/amazon-vpc-cni-k8s/pull/165

## Prerequisites

Before patching the node with the annotation, you will need to change the `AWS_VPC_K8S_CNI_CUSTOM_NETWORK_CFG` environment variable in the AWS CNI daemonset to `true`. By default, pods share the same subnet and security groups as the worker node's primary interface. When you set this variable to **true** it causes ipamD to use the security groups and VPC subnet in a worker node's ENIConfig for elastic network interface allocation. The subnet in the ENIConfig **must** belong to the same Availability Zone that the worker node resides in.

This operator can be configured to use an arbitrary EC2 tag for annotation.  If no tag is set, the operator will use a tag called k8s.amazonaws.com/eniConfig. While the operator will automatically apply the annotation to your worker nodes, it does not verify that subnet in the eniconfig is in the same AZ as the worker node.

The worker nodes also have to be assigned an IAM role with a policy that allows the DescribeTags action to the EC2 instance.  This is configured by default when you use eksctl to provision a cluster.

## Deploying

If you use `helm` to deploy applications you can use the
`charts/eniconfig-controller` chart to deploy the controller.

```bash
helm  install \
      --set controller.eniConfigName=<name-of-eniconfig-cr> \
      --set controller.automaticENIConfig=true \
      # --set controller.eniConfigName=default-eniconfig
      --set controller.eniConfigTagName=k8s.amazonaws.com/eniConfig
      https://github.com/christopherhein/eniconfig-controller/raw/helm-chart/eniconfig-controller-0.0.1.tgz
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