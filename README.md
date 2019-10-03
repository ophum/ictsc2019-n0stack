# ictsc2019-n0stack
testing...
VLANネットワークを作成した場合、そのネットワークのブリッジにホスト側のVLANインターフェースを接続する。
そのため、n0stack上のブリッジ内ではVLANタグはつかずホストから外に出るときにVLANタグが付与される。
デフォルトゲートウェイにはCIDRのうち最大のホストアドレスが利用される。


## install
`--vlan-external-interface`には、外部インターフェースを指定する。
外部インターフェースはL2スイッチにつながっており、Trunkとなる。
```
# go get -u github.com/ophum/ictsc2019-n0stack
# cd ~/go/bin/
# ./ictsc2019-n0stack install agent --arguments "--location=////1 --node-api-endpoint=localhost:20180 --vlan-external-interface=eth0"
```

## run
```
# n0core serve api & 
# systemctl start n0core-agent
```

## develop
`go/src/github.com/ophum/ictsc2019-n0stack/`にソースコードが配置されている。

### build
```
# cd go/src/github.com/ophum/ictsc2019-n0stack/
# go build
# ls ictsc2019-n0stack
```
