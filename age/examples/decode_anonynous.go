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
	yamlString := `
database_login: "service_1"
database_host: "db.company.com:5432"
database_password: !crypto/age |
  -----BEGIN AGE ENCRYPTED FILE-----
  YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHNjcnlwdCBvTDRrOUlXRGFYcXkzaVZu
  WXpzZndRIDE4ClZ3YVVHb0lVWlJtblVFazU4TlBkTitCWlg3dUNqd2N6R0hGVUFr
  T2gwb2sKLS0tIGFPYXBybWRUelNKeWkzc1lrVGpXUHJ4dDI4bWFDZEl6OXhpeTNY
  N0lIVjgKxPtRljkraTILjhf3v0MM5GmKnBwOMqLu2030RWMl6iW7YEYvunx2AMUA
  grTyTgUElzo=
  -----END AGE ENCRYPTED FILE-----`

	buf := bytes.NewBufferString(yamlString)

	id, err := age.NewScryptIdentity("point-adjust-member-tip-tiger-limb-honey-prefer-copy-issue")

	if err != nil {
		panic(err)
	}

	node := &yaml.Node{}

	w := ageyaml.Wrapper{
		Value:      node,
		Identities: []age.Identity{id},
	}

	decoder := yaml.NewDecoder(buf)
	err = decoder.Decode(&w)

	if err != nil {
		panic(err)
	}

	r, err := yaml.Marshal(&node)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", r)
}
