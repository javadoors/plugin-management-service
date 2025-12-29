/*
 * Copyright (c) 2024 Huawei Technologies Co., Ltd.
 * openFuyao is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"

	"plugin-management-service/pkg/client/k8s"
	"plugin-management-service/pkg/server/config"
	"plugin-management-service/pkg/server/runtime"
)

type dummyHandler struct{}

func (d dummyHandler) ServeHTTP(r http.ResponseWriter, req *http.Request) {
	// dummyHandler does nothing
}

func TestRecordAccessLogs(t *testing.T) {
	type args struct {
		req   *restful.Request
		resp  *restful.Response
		chain *restful.FilterChain
	}
	okFoundFilter := func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		resp.WriteHeader(http.StatusOK)
		chain.ProcessFilter(req, resp)
	}
	notFoundFilter := func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		resp.WriteHeader(http.StatusNotFound)
		chain.ProcessFilter(req, resp)
	}
	dummyTarget := func(req *restful.Request, resp *restful.Response) {
		// dummyTarget function does nothing
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"TestWarnLogs",
			args{
				restful.NewRequest(httptest.NewRequest("GET", "http://example.com", nil)),
				restful.NewResponse(httptest.NewRecorder()),
				&restful.FilterChain{
					Filters: []restful.FilterFunction{notFoundFilter},
					Target:  dummyTarget,
				},
			},
		},
		{
			"TestInfoLogs",
			args{
				restful.NewRequest(httptest.NewRequest("GET", "http://example.com", nil)),
				restful.NewResponse(httptest.NewRecorder()),
				&restful.FilterChain{
					Filters: []restful.FilterFunction{okFoundFilter},
					Target:  dummyTarget,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RecordAccessLogs(tt.args.req, tt.args.resp, tt.args.chain)
		})
	}
}

func TestInitServer(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.RunConfig
		wantErr bool
	}{
		{
			"TestInsecureConfig",
			&config.RunConfig{
				Server: &runtime.ServerConfig{
					BindAddress:    "0.0.0.0",
					SecurePort:     0,
					InsecurePort:   8080,
					PrivateKeyFile: "",
					CertFile:       "",
					CAFile:         "",
				},
				KubernetesCfg: k8s.NewKubernetesCfg(),
			},
			false,
		},
		{
			"TestSecureConfigNoFile",
			&config.RunConfig{
				Server: &runtime.ServerConfig{
					BindAddress:    "0.0.0.0",
					SecurePort:     8443,
					InsecurePort:   0,
					PrivateKeyFile: "",
					CertFile:       "",
					CAFile:         "",
				},
				KubernetesCfg: k8s.NewKubernetesCfg(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := initServer(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("initServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
