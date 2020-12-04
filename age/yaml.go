package age

import (
	"bytes"
	"io"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
	yamlv3 "gopkg.in/yaml.v3"
)

const (
	// YAMLTag ...
	YAMLTag = "!crypto/age"
)

// Wrapper is a struct that allows to decrypt crypted armored data in YAML as long
// that the data is tagged with `!crypto/age`.
//
//     agedata: !cripto/age |
//       -----BEGIN AGE ENCRYPTED FILE-----
//       ...
//       ...
//       -----END AGE ENCRYPTED FILE-----
//
type Wrapper struct {
	Value      interface{}
	Identities []age.Identity
}

// UnmarshalYAML ...
func (w Wrapper) UnmarshalYAML(value *yamlv3.Node) error {
	resolved, err := w.resolve(value)
	if err != nil {
		return err
	}

	return resolved.Decode(w.Value)
}

func (w Wrapper) resolve(node *yamlv3.Node) (*yamlv3.Node, error) {
	// Recurse into sequence types
	if node.Kind == yamlv3.SequenceNode || node.Kind == yamlv3.MappingNode {
		var err error

		if len(node.Content) > 0 {
			for i := range node.Content {
				node.Content[i], err = w.resolve(node.Content[i])
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if node.Tag != YAMLTag {
		return node, nil
	}

	// Check the absence of armored age header and footer
	valueTrimmed := strings.TrimSpace(node.Value)
	if !strings.HasPrefix(valueTrimmed, armor.Header) || !strings.HasSuffix(valueTrimmed, armor.Footer) {
		return node, nil
	}

	var armoredString string
	node.Decode(&armoredString)
	armoredStringReader := strings.NewReader(armoredString)
	armoredReader := armor.NewReader(armoredStringReader)
	decryptedReader, err := age.Decrypt(armoredReader, w.Identities...)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(decryptedReader)

	tempTag := node.Tag
	node.SetString(buf.String())
	node.Tag = tempTag

	return node, nil
}

// ArmoredString is a struct holding the string to crypt and the intended recipients
// of the encrypted output.
type ArmoredString struct {
	String     string
	Recipients []age.Recipient
}

// UnmarshalYAML ...
func (a *ArmoredString) UnmarshalYAML(value *yamlv3.Node) error {
	a.String = value.Value

	return nil
}

// MarshalYAML ...
func (a ArmoredString) MarshalYAML() (interface{}, error) {
	buf := &bytes.Buffer{}
	armorWriter := armor.NewWriter(buf)

	w, err := age.Encrypt(armorWriter, a.Recipients...)

	if err != nil {
		return nil, err
	}

	io.WriteString(w, string(a.String))
	w.Close()
	armorWriter.Close()

	node := yamlv3.Node{
		Kind:  yamlv3.ScalarNode,
		Tag:   YAMLTag,
		Value: string(buf.Bytes()),
	}

	return &node, nil
}

// EncryptYAML takes a Node and recursively marshal the Values.
func EncryptYAML(recipients []age.Recipient, node *yamlv3.Node) (*yamlv3.Node, error) {
	// Recurse into sequence types
	if node.Kind == yamlv3.SequenceNode || node.Kind == yamlv3.MappingNode {
		var err error

		if len(node.Content) > 0 {
			for i := range node.Content {
				node.Content[i], err = EncryptYAML(recipients, node.Content[i])
				if err != nil {
					return nil, err
				}
			}
		}

		return node, nil
	}

	if node.Tag != YAMLTag {
		return node, nil
	}

	armoredString := ArmoredString{String: node.Value, Recipients: recipients}
	nodeInterface, err := armoredString.MarshalYAML()

	return nodeInterface.(*yamlv3.Node), err
}
