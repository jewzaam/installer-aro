#!/bin/bash

# Remove agent-installer code
rm -rf ./cmd/openshift-install/
rm -rf ./pkg/agent
rm -rf ./pkg/asset/agent
rm -f ./pkg/asset/kubeconfig/agent.go
rm -rf ./data/data/agent
rm -rf ./pkg/asset/store/data/agent
rm -rf ./pkg/types/agent


# Remove agent-installer requires
sed -i "s\^replace github.com/openshift/assisted-service\//replace github.com/openshift/assisted-service\g" go.mod
sed -i "s\^	github.com/openshift/assisted-service\//    github.com/openshift/assisted-service\g" go.mod

# Remove tfvars
rm -f ./pkg/asset/cluster/tfvars.go
sed -i "s/&cluster.TerraformVariables{},//" ./pkg/asset/targets/targets.go

# Remove destroy (unused in ARO)
rm -rf ./pkg/destroy

K8S_VERSION=$(grep "replace k8s.io/client-go => k8s.io/client-go" go.mod | cut -d " " --fields="5")

# fix up the kube deps
go mod edit -replace k8s.io/kubectl=k8s.io/kubectl@$K8S_VERSION
go mod edit -replace k8s.io/api=k8s.io/api@$K8S_VERSION
go mod edit -replace k8s.io/apiserver=k8s.io/apiserver@$K8S_VERSION
go mod edit -replace k8s.io/apiextensions-apiserver=k8s.io/apiextensions-apiserver@$K8S_VERSION
go mod edit -replace k8s.io/component-base=k8s.io/component-base@$K8S_VERSION
go mod edit -replace k8s.io/apimachinery=k8s.io/apimachinery@$K8S_VERSION
go mod edit -replace k8s.io/code-generator=k8s.io/code-generator@$K8S_VERSION
go mod edit -replace k8s.io/kubelet=k8s.io/kubelet@$K8S_VERSION


go mod tidy
go mod vendor

cd ./hack/assets
go run ./
cd ../..
