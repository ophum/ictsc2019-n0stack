package iproute2

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

type Tap struct {
	name string
	queues int
	link netlink.Link // TODO
}

func NewTap(name string, queues uint32) (*Tap, error) {
	t := &Tap{
		name: name,
		queues: int(queues),
	}

	var err error
	t.link, err = netlink.LinkByName(name)
	if err != nil {
		if err := t.createTap(); err != nil {
			return nil, err
		}
	}

	if t.link.Type() != "tun" && t.link.Type() != "tuntap" {
		return nil, fmt.Errorf("Interface '%s' is not 'tun', but '%s'", name, t.link.Type())
	}

	return t, nil
}

func (t *Tap) createTap() error {
	l := netlink.NewLinkAttrs()
	l.Name = t.name
	t.link = &netlink.Tuntap{
		LinkAttrs: l,
		Mode:      netlink.TUNTAP_MODE_TAP,
		Flags:     netlink.TUNTAP_MULTI_QUEUE_DEFAULTS | netlink.TUNTAP_VNET_HDR,
		Queues:    t.queues,

	}

	if err := netlink.LinkAdd(t.link); err != nil {
		return fmt.Errorf("Failed 'ip tuntap add name %s mode tap': err='%s'", t.name, err.Error())
	}

	if err := t.Up(); err != nil {
		return err
	}

	return nil
}

func (t Tap) Name() string {
	return t.name
}

func (t *Tap) Up() error {
	if err := netlink.LinkSetUp(t.link); err != nil {
		return fmt.Errorf("Failed 'ip link set dev %s up': err='%s'", t.name, err.Error())
	}

	return nil
}

func (t *Tap) SetMaster(b *Bridge) error {
	if err := netlink.LinkSetMaster(t.link, b.link); err != nil {
		return fmt.Errorf("Failed ip link set dev %s master %s': err='%s'", t.name, b.name, err.Error())
	}

	return nil
}

func (t *Tap) Delete() error {
	if err := netlink.LinkDel(t.link); err != nil {
		return fmt.Errorf("Failed 'ip link del %s type bridge': err='%s'", t.name, err.Error())
	}

	t.link = nil
	return nil
}
