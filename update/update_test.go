package update

import (
	"reflect"
	"testing"
)

func Test_version_Newer(t *testing.T) {
	type args struct {
		o version
	}
	tests := []struct {
		name string
		v    version
		args args
		want bool
	}{
		{"newer", version{4, 18, 1}, args{version{4, 18, 0}}, true},
		{"eq", version{4, 18, 1}, args{version{4, 18, 1}}, false},
		{"old", version{4, 18, 1}, args{version{4, 18, 2}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Newer(&tt.args.o); got != tt.want {
				t.Errorf("version.Newer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strToVersion(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *version
		wantErr bool
	}{
		{"4.18.1", args{"4.18.1"}, &version{4, 18, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToVersion(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("strToVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("strToVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetLatestVersion(t *testing.T) {
	tests := []struct {
		name    string
		want    *version
		wantErr bool
	}{
		{"get", &version{4, 36, 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestVersion(repo + "tag/v4.36.2")
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatestVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLatestVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCore(t *testing.T) {
	type args struct {
		v *version
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"download", args{&version{4, 36, 2}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateCore(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
