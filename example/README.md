# example

- `image.yaml`: `/root/bionic-server-cloudimg-amd64.img`をイメージとして登録する
- `net.yaml`: networkを登録する。nativeとvlan100,200
- `vm.yaml`: vmを作成する。

```
n0cli do image.yaml
n0cli do net.yaml
n0cli do vm.yaml
```