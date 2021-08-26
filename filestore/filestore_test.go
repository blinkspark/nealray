package filestore

import (
	"testing"
)

func TestNewFSFileStore(t *testing.T) {
	type args struct {
		keyPath     string
		storePath   string
		listenAddrs []string
	}
	tests := []struct {
		name    string
		args    args
		want    *FSFileStore
		wantErr bool
	}{
		{name: "test", args: args{keyPath: "id.key", storePath: "data", listenAddrs: []string{}}, want: &FSFileStore{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFSFileStore(tt.args.keyPath, tt.args.storePath, tt.args.listenAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFSFileStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NewFSFileStore() got = %v, want %v", got, tt.want)
			}
		})
	}
}
