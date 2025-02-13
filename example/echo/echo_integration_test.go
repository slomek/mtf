// +build mtf

package framework

import (
	"testing"
	"time"

	"github.com/smallinsky/mtf/framework"
	"github.com/smallinsky/mtf/port"
	pb "github.com/smallinsky/mtf/proto/echo"
	pbo "github.com/smallinsky/mtf/proto/oracle"
)

func TestMain(m *testing.M) {
	framework.NewSuite(m).WithSut(framework.SutSettings{
		Dir:   "./service",
		Ports: []int{8001},
		Envs: []string{
			"ORACLE_ADDR=" + framework.GetDockerHostAddr(8002),
		},
	}).WithRedis(framework.RedisSettings{
		Password: "test",
	}).WithMySQL(framework.MysqlSettings{
		DatabaseName: "test_db",
		MigrationDir: "./service/migrations",
		Password:     "test",
	}).Run()
}

func TestEchoService(t *testing.T) {
	framework.Run(t, new(SuiteTest))
}

func (st *SuiteTest) Init(t *testing.T) {
	var err error
	if st.echoPort, err = port.NewGRPCClientPort((*pb.EchoClient)(nil), "localhost:8001"); err != nil {
		t.Fatalf("failed to init grpc client port")
	}
	if st.httpPort, err = port.NewHTTPPort(port.WithTLSHost("*.icndb.com")); err != nil {
		t.Fatalf("failed to init http port")
	}
	if st.oraclePort, err = port.NewGRPCServerPort((*pbo.OracleServer)(nil), ":8002"); err != nil {
		t.Fatalf("failed to init grpc oracle server")
	}
	time.Sleep(time.Millisecond * 300)
}

type SuiteTest struct {
	echoPort   *port.Port
	httpPort   *port.Port
	oraclePort *port.Port
}

func (st *SuiteTest) TestRedis(t *testing.T) {
	st.echoPort.Send(t, &pb.AskRedisRequest{
		Data: "make me sandwitch",
	})
	st.echoPort.Receive(t, &pb.AskRedisResponse{
		Data: "what? make it yourself",
	})
	st.echoPort.Send(t, &pb.AskRedisRequest{
		Data: "sudo make me sandwitch",
	})
	st.echoPort.Receive(t, &pb.AskRedisResponse{
		Data: "okey",
	})
}

func (st *SuiteTest) TestHTTP(t *testing.T) {
	st.echoPort.Send(t, &pb.AskGoogleRequest{
		Data: "Get answer for ultimate question of life the universe and everything",
	})
	st.httpPort.Receive(t, &port.HTTPRequest{
		Body:   []byte{},
		Method: "GET",
		Host:   "api.icndb.com",
		URL:    "/jokes/random?firstName=John\u0026amp;lastName=Doe",
	})
	st.httpPort.Send(t, &port.HTTPResponse{
		Body: []byte(`{"value":{"joke":"42"}}`),
	})
	st.echoPort.Receive(t, &pb.AskGoogleResponse{
		Data: "42",
	})
}

func (st *SuiteTest) TestClientServerGRPC(t *testing.T) {
	time.Sleep(time.Second * 2)
	st.echoPort.Send(t, &pb.AskOracleRequest{
		Data: "Get answer for ultimate question of life the universe and everything",
	})
	st.oraclePort.Receive(t, &pbo.AskDeepThroughRequest{
		Data: "Get answer for ultimate question of life the universe and everything",
	})
	st.oraclePort.Send(t, &pbo.AskDeepThroughRespnse{
		Data: "42",
	})
	st.echoPort.Receive(t, &pb.AskOracleResponse{
		Data: "42",
	})
}

func (st *SuiteTest) TestFetchDataFromDB(t *testing.T) {
	st.echoPort.Send(t, &pb.AskDBRequest{
		Data: "the dirty fork",
	})
	st.echoPort.Receive(t, &pb.AskDBResponse{
		Data: "Lucky we didn't say anything about the dirty knife",
	})
}
