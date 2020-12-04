module github.com/sylr/go-yaml-crypto/age

go 1.15

require (
	filippo.io/age v1.0.0-beta5
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace filippo.io/age => github.com/sylr/age v1.0.0-beta5.0.20201126225131-a495df083bec
