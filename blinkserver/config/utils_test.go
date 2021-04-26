package config

import (
	"reflect"
	"testing"
)

func TestGenDefaultConfig(t *testing.T) {
	GenDefaultConfig()
}

func TestLoadConfig(t *testing.T) {
	type args struct {
		fname string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{"read real", args{"config.json"}, &Config{[]ServerConfig{{Addr: ":80"}}}, false},
		{"read fake", args{"config_fake.json"}, &Config{}, false},
		{"read notexist", args{"config_notexist.json"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.fname)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
