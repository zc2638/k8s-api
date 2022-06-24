# CRD

根据 kubernetes 基础资源定义生成 CRD 资源

## 操作

1. 将缺失的 `+groupName` 注释补齐，core目录下的 `+groupName=core`

2. 安装 controller-gen 工具，参考 [https://github.com/kubernetes-sigs/controller-tools](https://github.com/kubernetes-sigs/controller-tools)
```shell
go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
```

4. 生成所有 CRD 资源
```shell
controller-gen crd paths=./... output:crd:artifacts:config=.output
```

4. 转换资源格式
```shell
go run cmd/convert-gen/main.go -src .output -dest .dest
```