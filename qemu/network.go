package qemu

import (
	"encoding/json"
	"fmt"
	"net"
)

// (QEMU) netdev_add id=tap0 type=tap vhost=true ifname=tap0 script=no downscript=no
// (QEMU) device_add driver=virtio-net-pci netdev=tap0 id=test0 mac=52:54:00:df:89:29 bus=pci.0
// まだべき等ではない
// TODO:
//   - すでにアタッチされていた場合、エラー処理を文字列で判定する必要がある
//   - MACアドレスを変更する
func (q Qemu) AttachTap(label, tap string, mac net.HardwareAddr, queues, vectors uint32) error {
	netdevID := fmt.Sprintf("netdev-%s", label)
	devID := fmt.Sprintf("virtio-net-%s", label)

	// check to create netdev

	if err := q.tapNetdevAdd(netdevID, tap, queues); err != nil {
		return fmt.Errorf("Failed to run netdev_add: err='%s'", err.Error())
	}

	if err := q.virtioNetPCIAdd(devID, netdevID, mac, vectors); err != nil {
		return fmt.Errorf("Failed to create virtio network device: err='%s'", err.Error())
	}

	return nil
}

func (q *Qemu) tapNetdevAdd(id, ifname string, queues uint32) error {
	cmd := struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Ifname     string `json:"ifname"`
		Vhost      bool   `json:"vhost"`
		Script     string `json:"script"`
		Downscript string `json:"downscript"`
		Queues     string `json:"queues"`
	}{
		id,
		"tap",
		ifname,
		true,
		"no",
		"no",
		fmt.Sprintf("%d", queues),
	}
	bs, err := json.Marshal(map[string]interface{}{
		"execute":   "netdev_add",
		"arguments": cmd,
	})

	if err != nil {
		return err
	}

	_, err = q.qmp.Run(bs)
	if err != nil {
		return err
	}

	return err
}

func (q *Qemu) virtioNetPCIAdd(devID, netdevID string, mac net.HardwareAddr, vectors uint32) error {
	cmd := struct {
		Driver  string `json:"driver"`
		ID      string `json:"id"`
		Netdev  string `json:"netdev"`
		Bus     string `json:"bus"`
		Mac     string `json:"mac"`
		Vectors string `json:"vectors"`
		Mq      string `json:"mq"`
		Guest_tso4 string `json:"guest_tso4"`
		Guest_tso6 string `json:"guest_tso4"`
		Guest_ecn  string `json:"guest_ecn"`
		Guest_ufo  string `json:"guest_ufo"`
	}{
		"virtio-net-pci",
		devID,
		netdevID,
		"pci.0",
		mac.String(),
		fmt.Sprintf("%d", vectors),
		"on",
		"off",
		"off",
		"off",
		"off",
	}

	bs, err := json.Marshal(map[string]interface{}{
		"execute":   "device_add",
		"arguments": cmd,
	})
	if err != nil {
		return err
	}

	_, err = q.qmp.Run(bs)
	if err != nil {
		return err
	}

	return nil
}
