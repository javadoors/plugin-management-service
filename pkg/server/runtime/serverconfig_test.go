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

package runtime

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		want *ServerConfig
	}{
		{
			name: "TestNewServer",
			want: NewServerConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServerConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerConfigValidate(t *testing.T) {
	_, errCert := os.Stat("/ssl/server.crt")
	if errCert == nil {
		t.Fatal("cert /ssl/server.crt should not exists")
	}
	_, errKey := os.Stat("/ssl/server.key")
	if errKey == nil {
		t.Fatal("cert key /ssl/key.crt should not exists")
	}

	tests := []struct {
		name         string
		BindAddress  string
		InsecurePort int
		SecurePort   int
		CertFile     string
		PrivateKey   string
		want         []error
	}{
		{
			name:         "TestValidateZeroPorts",
			BindAddress:  "0.0.0.0",
			InsecurePort: 0,
			SecurePort:   0,
			CertFile:     "",
			PrivateKey:   "",
			want: []error{
				fmt.Errorf("insecure and secure port can not be disabled at the same time"),
			},
		},
		{
			name:         "TestValidateInsecure",
			BindAddress:  "0.0.0.0",
			InsecurePort: 9032,
			SecurePort:   0,
			CertFile:     "",
			PrivateKey:   "",
			want:         []error{},
		},
		{
			name:         "TestValidateSecureEmptyCertNoKey",
			BindAddress:  "0.0.0.0",
			InsecurePort: 0,
			SecurePort:   4443,
			CertFile:     "",
			PrivateKey:   "/ssl/server.key",
			want: []error{
				fmt.Errorf("tls certificate file is empty while secure serving"),
				errKey,
			},
		},
		{
			name:         "TestValidateSecureNoCertEmptyKey",
			BindAddress:  "0.0.0.0",
			InsecurePort: 0,
			SecurePort:   4443,
			CertFile:     "/ssl/server.crt",
			PrivateKey:   "",
			want: []error{
				errCert,
				fmt.Errorf("tls private key file is empty while secure serving"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerConfig{
				BindAddress:    tt.BindAddress,
				InsecurePort:   tt.InsecurePort,
				SecurePort:     tt.SecurePort,
				CertFile:       tt.CertFile,
				PrivateKeyFile: tt.PrivateKey,
			}
			if got := s.Validate(); len(got) != 0 {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Validate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
