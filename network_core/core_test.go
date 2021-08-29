package networkcore

import (
	"testing"
)

func TestNewNetworkCore(t *testing.T) {
	type args struct {
		keyPath     string
		listenAddrs []string
	}
	tests := []struct {
		name    string
		args    args
		want    *NetworkCore
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", args{"id.key", []string{"/ip4/0.0.0.0/tcp/22333", "/ip6/::/tcp/22334"}}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNetworkCore(tt.args.keyPath, tt.args.listenAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNetworkCore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NewNetworkCore() = %v, want %v", got, tt.want)
			}
		})
	}
}
