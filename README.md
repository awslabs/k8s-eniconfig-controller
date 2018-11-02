# ENIConfig Controller

> This controller is still in an alpha state, please file issues and pull
> requests as you run into issues. Thanks 🎉

This repository will inplement auto annotating your Kubernete data plane nodes
with a desired `ENIConfig` name. This was originally implemented in the
`amazon-vpc-cni-k8s` project in this pull request -
https://github.com/aws/amazon-vpc-cni-k8s/pull/165

## Prerequisites

Before patching the node with the annotation, you will need to change the
`AWS_VPC_K8S_CNI_CUSTOM_NETWORK_CFG` environment variable in the AWS CNI
daemonset to `true`. By default, pods share the same subnet and security groups
as the worker node's primary interface. When you set this variable to **true**
it causes ipamD to use the security groups and VPC subnet in a worker node's
ENIConfig for elastic network interface allocation. The subnet in the ENIConfig
**must** belong to the same Availability Zone that the worker node resides in.

This operator can be configured to use an arbitrary EC2 tag for annotation.  If
no tag is set, the operator will use a tag called `k8s.amazonaws.com/eniConfig`.
While the operator will automatically apply the annotation to your worker nodes,
it does not verify that subnet in the eniconfig is in the same AZ as the worker
node.

The worker nodes also have to be assigned an IAM role with a policy that allows
the DescribeTags action to the EC2 instance.  This is configured by default when
you use eksctl to provision a cluster.

## Deploying

Before you deploy the controller, you will need to make sure your instances have
the proper policies assigned.

```bash
POLICY_ARN=$(aws iam create-policy \
                 --policy-name eniconfig-controller-policy \
                 --cli-input-json file://configs/eniconfig-controller-policy.json | jq -r ".Policy.Arn")
```

Now that you have this defined you can add this to the worker node role,
alternatively you could use a pod identity project to allow the pod to assume
your role.

```bash
aws iam attach-role-policy \
    --role-arn {WORKER NODE ROLE ARN} \
    --policy-arn $POLICY_ARN
```

Now that your permissions are properly configured you can deploy the controller
either using `helm` or manually by `kubectl`

If you use `helm` to deploy applications you can use the
`charts/eniconfig-controller` chart to deploy the controller.

```bash
helm install --name eniconfig-controller ./charts/eniconfig-controller
```

> Once this project is out of alpha state these will be available via the
> standard `helm` repositories.

If you'd like to use this without `helm` you can just `kubectl` apply the
manifest from the repository.

```bash
kubectl apply -f https://raw.githubusercontent.com/christopherhein/eniconfig-controller/master/configs/eniconfig-controller.yaml
```

## Running in Dev

```bash
# assumes you have a working kubeconfig, not required if operating in-cluster
$ go build -o eniconfig-controller .
$ ./eniconfig-controller -kubeconfig=$KUBECONFIG
```

## Releasing

To release this project all you have to do is run `make release`.

## License

This library is licensed under the Apache 2.0 License.
