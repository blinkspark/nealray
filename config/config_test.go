package config

import "testing"

func TestInitConfig(t *testing.T) {
	t.Run("TestInitConfig", func(t *testing.T) {
		if err := InitConfig(); err != nil {
			t.Errorf("InitConfig() error = %v, wantErr %v", err, nil)
		}
	})
}
