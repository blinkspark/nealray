package config

type Config struct {
	Servers []ServerConfig `json:"servers,omitempty"`
}

type ServerConfig struct {
	Addr string    `json:"addr,omitempty"`
	Root string    `json:"root,omitempty"`
	To   string    `json:"to,omitempty"`
	Cert string    `json:"cert,omitempty"`
	Key  string    `json:"key,omitempty"`
	DAV  DavConfig `json:"dav,omitempty"`
}

type DavConfig struct {
	Enable        bool      `json:"enable"`
	AccessControl bool      `json:"access_control,omitempty"`
	Users         []DavUser `json:"users,omitempty"`
}

type DavUser struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}
