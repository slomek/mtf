package framework

import (
	"testing"

	pb "github.com/smallinsky/mtf/e2e/proto/echo"
	"github.com/smallinsky/mtf/port"
)

func TestMain(m *testing.M) {
	NewSuite("suite_first", m).Run()
}

type SuiteTest struct {
	echoPort *port.ClientPort
	httpPort *port.HTTPPort
}

func (st *SuiteTest) Init(t *testing.T) {
	var err error
	if st.echoPort, err = port.NewGRPCClient((*pb.EchoClient)(nil), "localhost:8001"); err != nil {
		t.Fatalf("failed to init grpc client port")
	}
	if st.httpPort, err = port.NewHTTP(); err != nil {
		t.Fatalf("failed to init http port")
	}
}

func (st *SuiteTest) TestRedis(t *testing.T) {
	st.echoPort.SendT(t, &pb.AskRedisRequest{
		Data: "make me sandwitch",
	})

	st.echoPort.ReceiveT(t, &pb.AskRedisResponse{
		Data: "what? make it yourself",
	})

	st.echoPort.SendT(t, &pb.AskRedisRequest{
		Data: "sudo make me sandwitch",
	})
	st.echoPort.ReceiveT(t, &pb.AskRedisResponse{
		Data: "okey",
	})
}

func (st *SuiteTest) TestHTTP(t *testing.T) {
	st.echoPort.SendT(t, &pb.AskGoogleRequest{
		Data: "Get answer for ultimate question of life the universe and everything",
	})
	st.httpPort.ReceiveT(t, &port.HTTPRequest{
		Method: "GET",
	})
	st.httpPort.SendT(t, &port.HTTPResponse{
		Body: []byte(`{"value":{"joke":"42"}}`),
	})
	st.echoPort.ReceiveT(t, &pb.AskGoogleResponse{
		Data: "42",
	})
}

func TestEchoService(t *testing.T) {
	Run(t, new(SuiteTest))
}
