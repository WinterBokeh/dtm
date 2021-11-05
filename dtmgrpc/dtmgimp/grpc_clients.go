package dtmgimp

import (
	sync "sync"

	"github.com/yedf/dtm/dtmcli/dtmimp"
	grpc "google.golang.org/grpc"
)

type rawCodec struct{}

func (cb rawCodec) Marshal(v interface{}) ([]byte, error) {
	return v.([]byte), nil
}

func (cb rawCodec) Unmarshal(data []byte, v interface{}) error {
	ba, _ := v.([]byte)
	for index, byte := range data {
		ba[index] = byte
	}
	return nil
}

func (cb rawCodec) Name() string { return "dtm_raw" }

var normalClients, rawClients sync.Map

// MustGetDtmClient 1
func MustGetDtmClient(grpcServer string) DtmClient {
	return NewDtmClient(MustGetGrpcConn(grpcServer, false))
}

// MustGetRawDtmClient must get raw codec grpc conn
func MustGetRawDtmClient(grpcServer string) DtmClient {
	return NewDtmClient(MustGetGrpcConn(grpcServer, true))
}

// GetGrpcConn 1
func GetGrpcConn(grpcServer string, isRaw bool) (conn *grpc.ClientConn, rerr error) {
	clients := &normalClients
	if isRaw {
		clients = &rawClients
	}
	grpcServer = dtmimp.MayReplaceLocalhost(grpcServer)
	v, ok := clients.Load(grpcServer)
	if !ok {
		opts := grpc.WithDefaultCallOptions()
		if isRaw {
			opts = grpc.WithDefaultCallOptions(grpc.ForceCodec(rawCodec{}))
		}
		dtmimp.Logf("grpc client connecting %s", grpcServer)
		conn, rerr := grpc.Dial(grpcServer, grpc.WithInsecure(), grpc.WithUnaryInterceptor(GrpcClientLog), opts)
		if rerr == nil {
			clients.Store(grpcServer, conn)
			v = conn
			dtmimp.Logf("grpc client inited for %s", grpcServer)
		}
	}
	return v.(*grpc.ClientConn), rerr
}

// MustGetGrpcConn 1
func MustGetGrpcConn(grpcServer string, isRaw bool) *grpc.ClientConn {
	conn, err := GetGrpcConn(grpcServer, isRaw)
	dtmimp.E2P(err)
	return conn
}
