package config

type Config struct {
	Servers []ServerConfig `json:"servers,omitempty"`
}

type ServerConfig struct {
	Addr           string      `json:"addr,omitempty"`
	Root           string      `json:"root,omitempty"`
	EnableGzip     bool        `json:"enable_gzip,omitempty"`
	DisableListing bool        `json:"disable_listing,omitempty"`
	To             string      `json:"to,omitempty"`
	Cert           string      `json:"cert,omitempty"`
	Key            string      `json:"key,omitempty"`
	WebDAVs        []DavConfig `json:"webdavs,omitempty"`
}

type DavConfig struct {
	Root   string    `json:"root,omitempty"`
	Prefix string    `json:"prefix,omitempty"`
	Users  []DavUser `json:"users,omitempty"`
}

type DavUser struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}
