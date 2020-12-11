package age

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"filippo.io/age"
	"filippo.io/age/armor"
	"gopkg.in/yaml.v3"
)

func TestSimpleDataString(t *testing.T) {
	// Key file
	keyFile, err := os.Open("./testdata/age.key")

	if err != nil {
		t.Fatal(err)
	}

	// Parse key files for identities
	ids, err := age.ParseIdentities(keyFile)

	if err != nil {
		t.Fatal(err)
	}

	// Parse key files for recipients
	keyFile.Seek(0, io.SeekStart)
	recs, err := age.ParseRecipients(keyFile)

	if err != nil {
		t.Fatal(err)
	}

	d1 := struct {
		Data ArmoredString `yaml:"data"`
	}{
		Data: ArmoredString{
			Value:      "this is a test",
			Recipients: recs,
		},
	}

	// Marshal
	bytes, err := yaml.Marshal(&d1)

	if err != nil {
		t.Fatal(err)
	}

	str := string(bytes)

	if strings.Index(str, armor.Header) == -1 || strings.Index(str, armor.Footer) == -1 {
		t.Errorf("Armored Age header or footer missing in yaml:\n%s", str)
		t.FailNow()
	}

	// Unmarshal
	d2 := struct {
		Data string `yaml:"data"`
	}{}

	w := Wrapper{
		Value:      &d2,
		Identities: ids,
	}

	err = yaml.Unmarshal(bytes, &w)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	Convey("Compare orginal struct with result of Unmarshalling", t, func() {
		So(d1.Data.String(), ShouldEqual, d2.Data)
	})
}

func TestSimpleDataArmoredString(t *testing.T) {
	// Key file
	keyFile, err := os.Open("./testdata/age.key")

	if err != nil {
		t.Fatal(err)
	}

	// Parse key files for identities
	ids, err := age.ParseIdentities(keyFile)

	if err != nil {
		t.Fatal(err)
	}

	// Parse key files for recipients
	keyFile.Seek(0, io.SeekStart)
	recs, err := age.ParseRecipients(keyFile)

	if err != nil {
		t.Fatal(err)
	}

	d1 := struct {
		Data ArmoredString `yaml:"data"`
	}{
		Data: ArmoredString{
			Recipients: recs,
			Value:      "this is a test",
		},
	}

	// Marshal
	bytes, err := yaml.Marshal(&d1)

	if err != nil {
		t.Fatal(err)
	}

	str := string(bytes)

	Convey("Check age armor header and footer", t, FailureHalts, func() {
		So(str, ShouldContainSubstring, armor.Header)
		So(str, ShouldContainSubstring, armor.Footer)
	})

	// Unmarshal
	d2 := struct {
		Data string `yaml:"data"`
	}{}

	w := Wrapper{
		Value:      &d2,
		Identities: ids,
	}

	err = yaml.Unmarshal(bytes, &w)

	Convey("Compare orginal struct with result of Unmarshalling", t, func() {
		So(err, ShouldBeNil)
		So(d1.Data.String(), ShouldEqual, d2.Data)
	})
}

func TestAnonymousStruct(t *testing.T) {
	// Add identities
	keyFile, err := os.Open("./testdata/age.key")

	if err != nil {
		t.Fatal(err)
	}

	ids, err := age.ParseIdentities(keyFile)

	// Open source yaml
	yamlFile, err := os.Open("./testdata/lipsum.yaml")

	if err != nil {
		t.Fatal(err)
	}

	// "anonymous" struct
	d1 := make(map[interface{}]interface{})

	// Decode
	w := Wrapper{Value: &d1, Identities: ids}
	decoder := yaml.NewDecoder(yamlFile)
	decoder.Decode(&w)

	// Check that the decoded yaml has the lipsum key
	if _, ok := d1["lipsum"]; !ok {
		t.Errorf("Decoded yaml has no lipsum key")
		t.FailNow()
	}

	// Open the file containing the original data used in the yaml file
	lipsumFile, err := os.Open("./testdata/lipsum.txt")

	if err != nil {
		t.Fatal(err)
	}

	lipsumBuf, err := ioutil.ReadAll(lipsumFile)

	if err != nil {
		t.Fatal(err)
	}

	lipsum := string(lipsumBuf)

	Convey("Compare orginal lipsum to decoded one", t, func() {
		So(d1["lipsum"], ShouldEqual, lipsum)
	})
}

type complexStruct struct {
	RegularData []string        `yaml:"regularData"`
	CryptedData []ArmoredString `yaml:"cryptedData"`
}

func TestComplexData(t *testing.T) {
	keyFile, err := os.Open("./testdata/age.key")

	if err != nil {
		t.Fatal(err)
	}

	ids, err := age.ParseIdentities(keyFile)

	if err != nil {
		t.Fatal(err)
	}

	keyFile.Seek(0, io.SeekStart)
	recs, err := age.ParseRecipients(keyFile)

	if err != nil {
		t.Fatal(err)
	}

	// -- test 1 ---------------------------------------------------------------

	d1 := complexStruct{
		RegularData: []string{
			"this is the first pwet",
			"this is the second pwet",
		},
		CryptedData: []ArmoredString{
			{Value: "this is supposed to be crypted", Recipients: recs},
			{Value: "this is also supposed to be crypted", Recipients: recs},
		},
	}

	out1, err := yaml.Marshal(&d1)

	Convey("Unmarshal should not return error", t, FailureHalts, func() {
		So(err, ShouldBeNil)
	})

	Convey("Search for non encrypted data which shouldn't be", t, func() {
		So(err, ShouldBeNil)
		So(string(out1), ShouldContainSubstring, "this is the first pwet")
		So(string(out1), ShouldContainSubstring, "this is the second pwet")
	})

	Convey("Search for non encrypted data which should be encrypted", t, func() {
		So(err, ShouldBeNil)
		So(string(out1), ShouldNotContainSubstring, "this is supposed to be crypted")
		So(string(out1), ShouldNotContainSubstring, "this is also supposed to be crypted")
	})

	// -- test 2 ---------------------------------------------------------------

	d2 := yaml.Node{}
	w := Wrapper{Value: &d2, Identities: ids}
	err = yaml.Unmarshal(out1, &w)

	Convey("Unmarshal should not return error", t, FailureHalts, func() {
		So(err, ShouldBeNil)
	})

	Convey("Search for encrypted data which shouldn't be", t, func() {
		var recurse func(node *yaml.Node)
		recurse = func(node *yaml.Node) {
			if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
				if len(node.Content) > 0 {
					for i := range node.Content {
						recurse(node.Content[i])
					}
				}
			}

			So(node.Value, ShouldNotContainSubstring, armor.Header)
			So(node.Value, ShouldNotContainSubstring, armor.Footer)
		}
		recurse(&d2)
	})

	// -- test 3 ---------------------------------------------------------------

	d3, err := MarshalYAML(&d2, recs)

	Convey("MarshalYAML should not return error", t, FailureHalts, func() {
		So(err, ShouldBeNil)
	})

	out2, err := yaml.Marshal(&d3)

	Convey("Marshalling should not return error", t, FailureHalts, func() {
		So(err, ShouldBeNil)
	})

	Convey("Search for non encrypted data which shouldn't be", t, func() {
		So(err, ShouldBeNil)
		So(string(out2), ShouldContainSubstring, "this is the first pwet")
		So(string(out2), ShouldContainSubstring, "this is the second pwet")
	})

	Convey("Search for non encrypted data which should be encrypted", t, func() {
		So(string(out2), ShouldNotContainSubstring, "this is supposed to be crypted")
		So(string(out2), ShouldNotContainSubstring, "this is also supposed to be crypted")
	})

	Convey("Compare orginal yaml to re-marshalled one, it should differ due to age rekeying", t, func() {
		So(string(out1), ShouldNotEqual, string(out2))
	})

	// -- test 4 ---------------------------------------------------------------

	d4 := complexStruct{}
	w.Value = &d4
	err = yaml.Unmarshal(out1, &w)

	Convey("Unmarshalling should not return error", t, FailureHalts, func() {
		So(err, ShouldBeNil)
	})

	Convey("Search for non encrypted data which should be", t, func() {
		So(d4.RegularData[0], ShouldContainSubstring, "this is the first pwet")
		So(d4.RegularData[1], ShouldContainSubstring, "this is the second pwet")
	})

	Convey("Search for encrypted data which shouldn't be", t, func() {
		So(d4.CryptedData[0].String(), ShouldContainSubstring, "this is supposed to be crypted")
		So(d4.CryptedData[1].String(), ShouldContainSubstring, "this is also supposed to be crypted")
	})
}

func TestUnlmarshallingBogusEncryptedData(t *testing.T) {
	tests := []struct {
		Description string
		Assertion   func(interface{}, ...interface{}) string
		YAML        string
	}{
		{
			Description: "Bogus age payload: bogus base64",
			Assertion:   ShouldBeError,
			YAML: `
database_login: "service_1"
database_host: "db.company.com:5432"
database_password: !crypto/age |
  -----BEGIN AGE ENCRYPTED FILE-----
  YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHNjcnlwdCBvTDRrOUlXRGFYcXkzaVZu
  WXpzZndRzIDE4ClZ3YVVHb0lVWlJtblVFazU4TlBkTitCWlg3dUNqd2N6R0hGVUFr
  T2gwb2sKLS0tIGFPYXBybWRUelNKeWkzc1lrVGpXUHJ4dDI4bWFDZEl6OXhpeTNY
  N0lIVjgKxPtRljkraTILjhf3v0MM5GmKnBwOMqLu2030RWMl6iW7YEYvunx2AMUA
  grTyTgUElzo=....
  -----END AGE ENCRYPTED FILE-----
`,
		},
		{
			Description: "Bogus age payload: no base64",
			Assertion:   ShouldBeError,
			YAML: `
database_login: "service_1"
database_host: "db.company.com:5432"
database_password: !crypto/age |
  -----BEGIN AGE ENCRYPTED FILE-----
  ...
  -----END AGE ENCRYPTED FILE-----
`,
		},
		{
			Description: "Bogus age payload: base64 not age data",
			Assertion:   ShouldBeError,
			YAML: `
database_login: "service_1"
database_host: "db.company.com:5432"
database_password: !crypto/age |
  -----BEGIN AGE ENCRYPTED FILE-----
  cWtsc2RobGtxZGhqc2ts
  -----END AGE ENCRYPTED FILE-----
`,
		},
		{
			Description: "Not encrypted payload",
			Assertion:   ShouldBeNil,
			YAML: `
database_login: "service_1"
database_host: "db.company.com:5432"
database_password: !crypto/age |
  this is a test
`,
		},
	}

	id, err := age.NewScryptIdentity("point-adjust-member-tip-tiger-limb-honey-prefer-copy-issue")

	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		buf := bytes.NewBufferString(test.YAML)
		node := yaml.Node{}

		w := Wrapper{
			Value:      &node,
			Identities: []age.Identity{id},
		}
		decoder := yaml.NewDecoder(buf)
		err = decoder.Decode(&w)

		Convey(test.Description, t, func() {
			So(err, test.Assertion)
		})
	}
}
