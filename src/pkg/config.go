package pkg

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Authn struct {
	Users []User `yaml:"users"`
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	OrgID    string `yaml:"orgid"`
}

func ParseConfig(location *string) (*Authn, error) {
	data, err := ioutil.ReadFile(*location)
	if err != nil {
		return nil, err
	}
	authn := Authn{}
	err = yaml.Unmarshal([]byte(data), &authn)
	if err != nil {
		return nil, err
	}
	return &authn, nil
}
