package ping

import (
	"context"
	ping_pb "easycoding/api/ping"
	"testing"
)

func TestPing(t *testing.T) {
	s := &service{}
	resp, err := s.Ping(context.Background(), &ping_pb.PingRequest{
		Req: "",
	})
	if err != nil {
		t.Errorf("Ping error %s", err.Error())
	}
	respMsg := "pong"
	if resp.Res != respMsg {
		t.Errorf("Response error, expect %s, got %s", respMsg, resp.Res)
	}
}
