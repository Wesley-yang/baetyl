// +build !windows

package master

import (
	"os"
	"testing"
	"time"

	"github.com/baidu/openedge/module/config"
	"github.com/baidu/openedge/module/master"
	"github.com/stretchr/testify/assert"
)

func TestAPIUnix(t *testing.T) {
	os.MkdirAll("./var/", 0755)
	defer os.RemoveAll("./var/")
	addr := "unix://./var/test.sock"
	s, err := NewServer(&mockAPI{pass: true}, config.HTTPServer{Address: addr, Timeout: time.Minute})
	assert.NoError(t, err)
	defer s.Close()
	err = s.Start()
	assert.NoError(t, err)
	c, err := master.NewClient(config.HTTPClient{Address: addr, Timeout: time.Minute, KeepAlive: time.Minute})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	p, err := c.GetPortAvailable("127.0.0.1")
	assert.NoError(t, err)
	assert.NotZero(t, p)
	err = c.StartModule(&config.Module{Name: "name"})
	assert.NoError(t, err)
	err = c.StopModule(&config.Module{Name: "name"})
	assert.NoError(t, err)
	stats, err := c.Stats()
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	err = c.Reload("var/db/v1.zip")
	assert.NoError(t, err)
}

func TestAPIUnixUnauthorized(t *testing.T) {
	os.MkdirAll("./var/", 0755)
	defer os.RemoveAll("./var/")
	addr := "unix://./var/test.sock"
	s, err := NewServer(&mockAPI{pass: false}, config.HTTPServer{Address: addr, Timeout: time.Minute})
	assert.NoError(t, err)
	defer s.Close()
	err = s.Start()
	assert.NoError(t, err)
	c, err := master.NewClient(config.HTTPClient{Address: addr, Timeout: time.Minute, KeepAlive: time.Minute, Username: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	_, err = c.GetPortAvailable("127.0.0.1")
	assert.EqualError(t, err, "[400] account (test) unauthorized")
	err = c.StartModule(&config.Module{Name: "name"})
	assert.EqualError(t, err, "[400] account (test) unauthorized")
	err = c.StopModule(&config.Module{Name: "name"})
	assert.EqualError(t, err, "[400] account (test) unauthorized")
	stats, err := c.Stats()
	assert.EqualError(t, err, "[400] account (test) unauthorized")
	assert.Nil(t, stats)
	err = c.Reload("var/db/v1.zip")
	assert.EqualError(t, err, "[400] account (test) unauthorized")
}
