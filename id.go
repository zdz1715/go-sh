package sh

import "github.com/rs/xid"

type IDCreator func() string

var XidCreator IDCreator = func() string {
	return xid.New().String()
}
