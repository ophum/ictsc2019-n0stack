# example

`request_node_name`にはnodeの名前を指定する。`n0cli get node`

- `image.yaml`: `/root/bionic-server-cloudimg-amd64.img`をイメージとして登録する
- `net.yaml`: networkを登録する。nativeとvlan100,200
- `vm.yaml`: vmを作成する。

```
n0cli do image.yaml
n0cli do net.yaml
n0cli do vm.yaml
```

## net-native
VLANではないネイティブなネットワーク。
`net-native-{uniqid}`なブリッジにアドレスを振ると、このネットワークに接続するVMにホストからアクセスすることができる。

## net-100
VLAN ID `100`のネットワーク。
`net-100-{uniqid}`なブリッジにVLAN ID `100`の通信で、このネットワークに接続するVMにアクセスできる。
また、`vlan-external-interface`で指定したインターフェースに、VLAN100のサブインターフェースが作成され、ブリッジに接続される。

## net-200
VLAN ID `200`のネットワーク。
`net-200-{uniqid}`なブリッジにVLAN ID `200`の通信で、このネットワークに接続するVMにアクセスできる。
また、`vlan-external-interface`で指定したインターフェースに、VLAN200のサブインターフェースが作成され、ブリッジに接続される。
