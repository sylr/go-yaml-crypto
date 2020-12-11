// +build ignore

package main

import (
	"bytes"
	"fmt"

	"filippo.io/age"
	ageyaml "github.com/sylr/go-yaml-crypto/age"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	rec, err := age.NewScryptRecipient("point-adjust-member-tip-tiger-limb-honey-prefer-copy-issue")

	if err != nil {
		panic(err)
	}

	node := struct {
		DatabaseLogin    string                `yaml:"database_login"`
		DatabaseHost     string                `yaml:"database_host"`
		DatabasePassword ageyaml.ArmoredString `yaml:"database_password"`
	}{
		DatabaseLogin: "service_1",
		DatabaseHost:  "db.company.com:5432",
		DatabasePassword: ageyaml.ArmoredString{
			Value:      "MyDatabasePassword",
			Recipients: []age.Recipient{rec},
		},
	}

	buf := bytes.NewBuffer(nil)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	err = encoder.Encode(&node)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", buf.String())
}
