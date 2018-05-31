package hiera

import (
	"encoding/json"
	"testing"

	"github.com/spf13/cast"
)

func TestHieraExec(t *testing.T) {
	var f interface{}

	hiera := testHieraConfig()

	out, err := hiera.Exec("aws_cloudwatch_enable")
	if err != nil {
		t.Errorf("Error running hiera: %s", err)
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %s", err)

	}

	v := cast.ToString(f)

	if v != "true" {
		t.Errorf("aws_cloudwatch_enable is %s; want %s", v, "true")
	}
}

func TestHieraHash(t *testing.T) {
	hiera := testHieraConfig()

	v, err := hiera.Hash("aws_tags")
	if err != nil {
		t.Errorf("Error running hiera.Hash: %s", err)
	}

	if v["team"] != "A" {
		t.Errorf("aws_tags.team is %s; want %s", v, "A")
	}

	if v["tier"] != "1" {
		t.Errorf("aws_tags.tier is %s; want %s", v, "1")
	}
}

func TestHieraValue(t *testing.T) {
	hiera := testHieraConfig()

	v, err := hiera.Value("aws_cloudwatch_enable")
	if err != nil {
		t.Errorf("Error running hiera.Value: %s", err)
	}

	if v != "true" {
		t.Errorf("aws_cloudwatch_enable is %s; want %s", v, "true")
	}
}

func testHieraConfig() Hiera {
	return NewHiera(
		"hiera",
		"../tests/hiera.yaml",
		map[string]interface{}{"service": "api", "environment": "live"},
	)
}
