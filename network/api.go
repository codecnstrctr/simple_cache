package network

import (
	"net/rpc"
	"time"

	"../cache"
)

const (
	okResponse = "OK"
)

// Command represents the handler
// of commands received from client
type Command struct {
	storage *cache.Cache
}

// SetRequest represents the request
// to set a value in cache
type SetRequest struct {
	Key  string
	Data interface{}
	TTL  time.Duration
}

// GetResponse represents the response
// on get cache request
type GetResponse struct {
	Key   string
	Found bool
	Data  interface{}
}

// NewCommand creates new command handler instance
func NewCommand(storage *cache.Cache) *Command {
	c := new(Command)
	c.storage = storage
	rpc.Register(c)

	return c
}

// Set is cache api set command
func (c *Command) Set(req *SetRequest, resp *string) error {
	c.storage.Set(req.Key, req.Data, req.TTL)
	*resp = okResponse

	return nil
}

// Get is cache api get command
func (c *Command) Get(req string, resp *GetResponse) error {
	resp.Data, resp.Found = c.storage.Get(req)
	resp.Key = req

	return nil
}

// Remove is cache api remove command
func (c *Command) Remove(req string, resp *bool) error {
	*resp = c.storage.Remove(req)

	return nil
}

// Keys is cache api get command
func (c *Command) Keys(req struct{}, resp *[]string) error {
	*resp = c.storage.Keys()

	return nil
}
