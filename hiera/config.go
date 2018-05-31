package hiera

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cast"
)

type Config struct {
	Bin    string
	Config string
	Scope  map[string]interface{}
}

func (c *Config) Exec(arg ...string) ([]byte, error) {
	var scope []string
	var args []string
	var out []byte
	var err error

	for key, value := range c.Scope {
		scope = append(scope, strings.Join([]string{key, value.(string)}, "="))
	}

	for _, r := range [][]string{[]string{"-f", "json", "-c", c.Config}, arg, scope} {
		args = append(args, r...)
	}

	if out, err = exec.Command(c.Bin, args...).Output(); err != nil {
		log.Println(err)
	}
	log.Printf("Out: %s\n", out)
	return out, nil
}

func (c *Config) Array(key string) ([]interface{}, error) {
	var f interface{}
	var e []interface{}

	out, err := c.Exec("-a", key)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for i, v := range f.([]interface{}) {
		log.Println(i)
		e = append(e, cast.ToString(v))
	}

	return e, nil
}

func (c *Config) Hash(key string) (map[string]interface{}, error) {
	var f interface{}

	out, err := c.Exec("-h", key)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(f.(map[string]interface{}))

	return f.(map[string]interface{}), nil
}

func (c *Config) Value(key string) (string, error) {
	var f interface{}

	out, err := c.Exec(key)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return cast.ToString(f), nil
}
