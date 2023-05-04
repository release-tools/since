/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cfg

import (
	"reflect"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    SinceConfig
		wantErr bool
	}{
		{
			name:    "no config file",
			args:    args{configPath: "testdata/no-config.yaml"},
			want:    SinceConfig{},
			wantErr: false,
		},
		{
			name:    "valid config file",
			args:    args{configPath: "testdata/valid-config.yaml"},
			want:    SinceConfig{Before: []Hook{{Command: "echo", Args: []string{"hello world"}}}},
			wantErr: false,
		},
		{
			name:    "invalid config file",
			args:    args{configPath: "testdata/invalid-config.yaml"},
			want:    SinceConfig{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfig(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
