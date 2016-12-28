package memds

import "github.com/BurntSushi/toml"

type Config struct {
	Port      int    `toml:"port"`
	Sock      string `toml:"sock"`
	BucketNum int    `toml:"bucket_num"`
}

func LoadConfig(p string) (*Config, error) {
	var c Config
	if _, err := toml.DecodeFile(p, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
