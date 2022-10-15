package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type readFromRequestFunc func(context.Context, *http.Request) []string

var funcs = []readFromRequestFunc{
	readAuthFromRequest,
}

func readFromRequest(logger *logrus.Logger) func(ctx context.Context, req *http.Request) metadata.MD {
	return func(ctx context.Context, req *http.Request) metadata.MD {
		data := []string{}
		for _, f := range funcs {
			pair := f(ctx, req)
			if len(pair)%2 == 0 {
				data = append(data, pair...)
			} else {
				logger.Warnf("error parsed request and header, got odd metadata %v\n", pair)
			}
		}
		return metadata.Pairs(data...)
	}
}

func readAuthFromRequest(ctx context.Context, req *http.Request) []string {
	authKey := "Authorization"
	authScheme := "Bearer"
	val := req.Header.Get(authKey)
	authInfo := []string{}
	if val == "" {
		return []string{}
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return []string{}
	}
	if !strings.EqualFold(splits[0], authScheme) {
		return []string{}
	}
	authInfo = append(authInfo, authKey, fmt.Sprintf("%s %s", authScheme, splits[1]))
	return authInfo
}
