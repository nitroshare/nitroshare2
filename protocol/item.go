package protocol

import (
	"io"
)

// Item is an interface for data being sent in a transfer. Any item being sent must implement this interface.
type Item interface {

	// Type is one of the ItemType* constants or a custom value in reverse domain name notation.
	Type() string

	// Size indicates the total size (in bytes) of the main content for this item.
	Size() int64

	// Meta is information that should be sent with the item header and which will be encoded as JSON data. For simple things like a URL, this can also contain the actual payload.
	Meta() (any, error)

	// Open returns a ReadCloser which will return the main content for the item. If there is no content, Size() should return 0 and this method may be a stub that returns nil, nil.
	Open() (io.ReadCloser, error)
}
