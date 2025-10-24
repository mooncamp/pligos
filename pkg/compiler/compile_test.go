package compiler

import (
	_ "embed"
	"testing"

	"realcloud.tech/pligos/pkg/maputil"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

//go:embed testdata/a/schema.yaml
var aSchemaYAML []byte

//go:embed testdata/a/pligos.yaml
var aPligosYAML []byte

//go:embed testdata/a/result.yaml
var aResultYAML []byte

//go:embed testdata/b/schema.yaml
var bSchemaYAML []byte

//go:embed testdata/b/pligos.yaml
var bPligosYAML []byte

//go:embed testdata/b/result.yaml
var bResultYAML []byte

func testCompile(schemaYaml, pligosYaml, resultYaml []byte, t *testing.T) {
	var schema map[string]interface{}
	if err := yaml.Unmarshal(schemaYaml, &schema); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(pligosYaml, &config); err != nil {
		t.Fatalf("unmarshal pligos config: %v", err)
	}

	var expected map[string]interface{}
	if err := yaml.Unmarshal(resultYaml, &expected); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	normalizer := &maputil.Normalizer{}
	c := &Compiler{
		config:    normalizer.Normalize(config)["contexts"].(map[string]interface{})["base"].(map[string]interface{}),
		schema:    normalizer.Normalize(schema)["context"].(map[string]interface{}),
		instances: normalizer.Normalize(config),
		types:     normalizer.Normalize(schema),
	}

	schema = normalizer.Normalize(schema)
	config = normalizer.Normalize(config)

	res, err := c.Compile()
	if err != nil {
		t.Fatalf("graph compile: %v", err)
	}

	assert.Equal(t, normalizer.Normalize(expected), res)
}

func Test_compile_a(t *testing.T) {
	testCompile(aSchemaYAML, aPligosYAML, aResultYAML, t)
}

func Test_compile_b(t *testing.T) {
	testCompile(bSchemaYAML, bPligosYAML, bResultYAML, t)
}
