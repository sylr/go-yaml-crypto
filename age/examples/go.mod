module github.com/sylr/go-yaml-crypto/age/examples

go 1.15

require (
	filippo.io/age v1.0.0-beta5
	github.com/sylr/go-yaml-crypto/age v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace (
	github.com/sylr/go-yaml-crypto/age => ../
	gopkg.in/yaml.v3 => github.com/sylr/go-yaml v0.0.0-20201211202443-be0157e6a8ed
)
