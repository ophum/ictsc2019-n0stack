package ictscnetwork

import (
	"errors"
	"context"

	"google.golang.org/grpc/codes"

	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/network"
	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	netutil "github.com/n0stack/n0stack/n0core/pkg/util/net"
	ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"
)

type NetworkAPI struct {
	*network.NetworkAPI
	dataStore datastore.Datastore
}

func CreateNetworkAPI(ds datastore.Datastore) *NetworkAPI {
	a := &NetworkAPI{
		NetworkAPI: network.CreateNetworkAPI(ds),
		dataStore: ds.AddPrefix("network"),
	}

	return a
}

func (a NetworkAPI) ApplyNetwork(ctx context.Context, req *ppool.ApplyNetworkRequest) (*ppool.Network, error) {
	ipv4 := netutil.ParseCIDR(req.Ipv4Cidr)
	ipv6 := netutil.ParseCIDR(req.Ipv6Cidr)
	{
		if req.Name == "" {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set any 'name'")
		}

		if ipv4 == nil && ipv6 == nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Failed 'ipv4_cidr' and 'ipv6_cidr' are invalid")
		}
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, errors.New("hoge")
	}
	defer a.dataStore.Unlock(req.Name)

	network := &ppool.Network{}
	if err := a.dataStore.Get(req.Name, network); err != nil  && !datastore.IsNotFound(err) {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get data from db: err='%s'", err.Error())
	}

	{
		if network.Name != "" && ipv4.String() != network.Ipv4Cidr {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "ipv4 cidr is different from previous ipv4 cidr")
		}
	}

	network.Name = req.Name
	network.Annotations = req.Annotations
	network.Ipv4Cidr = req.Ipv4Cidr
	network.Ipv6Cidr = req.Ipv6Cidr
	network.Domain = req.Domain

	network.State = ppool.Network_AVAILABLE
	if err := a.dataStore.Apply(req.Name, network); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}

	//return network, nil
	return nil, errors.New("hoge")
}
