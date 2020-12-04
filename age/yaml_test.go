package age

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/andreyvit/diff"
	"gopkg.in/yaml.v3"
	yamlv3 "gopkg.in/yaml.v3"
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
			String:     "this is a test",
			Recipients: recs,
		},
	}

	// Marshal
	bytes, err := yamlv3.Marshal(&d1)

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

	err = yamlv3.Unmarshal(bytes, &w)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if d1.Data.String != d2.Data {
		t.Errorf("Expected `%s` got `%s`", d1.Data, d2.Data)
		t.FailNow()
	}
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
			String:     "this is a test",
		},
	}

	// Marshal
	bytes, err := yamlv3.Marshal(&d1)

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

	err = yamlv3.Unmarshal(bytes, &w)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if d1.Data.String != d2.Data {
		t.Errorf("Expected `%s` got `%s`", d1.Data, d2.Data)
		t.FailNow()
	}
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
	decoder := yamlv3.NewDecoder(yamlFile)
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

	// Compare original and decrypted data
	if d1["lipsum"] != lipsum {
		t.Errorf("Original and decrypted lipsum differ: %s", diff.CharacterDiff(lipsum, fmt.Sprintf("%v", d1["lipsum"])))
		t.FailNow()
	}
}

type complexStruct struct {
	RegularData []string        `yaml:"regularData"`
	CryptedData []ArmoredString `yaml:"cryptedData"`
}

func TestGlobalComplexData(t *testing.T) {
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

	d1 := complexStruct{
		RegularData: []string{
			"this is the first pwet",
			"this is the second pwet",
		},
		CryptedData: []ArmoredString{
			{String: "this is supposed to be crypted", Recipients: recs},
			{String: "this is also supposed to be crypted", Recipients: recs},
		},
	}

	out, err := yamlv3.Marshal(&d1)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", string(out))

	d2 := yaml.Node{}

	w := Wrapper{Value: &d2, Identities: ids}
	err = yamlv3.Unmarshal(out, &w)

	if err != nil {
		t.Fatal(err)
	}

	d3, err := EncryptYAML(recs, &d2)

	if err != nil {
		t.Fatal(err)
	}

	out, err = yamlv3.Marshal(&d3)

	t.Logf("%s", string(out))
}
