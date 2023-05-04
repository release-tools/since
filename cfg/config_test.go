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
