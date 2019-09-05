// Package uuid is a utility package to standardize and abstract away how UUIDs are generated and used.
package uuid

import (
	gid "github.com/satori/go.uuid"
)

// UUID is a wrapper around go.uuid.UUID
type UUID struct {
	gid.UUID
}

// New generates a new UUID.
func New() UUID {
	id := gid.NewV4()
	return UUID{id}
}

// FromString reads a UUID from a string
func FromString(str string) (UUID, error) {
	id, err := gid.FromString(str)
	if err != nil {
		return UUID{}, err
	}
	return UUID{id}, nil
}

// ToString converts UUID into string
func (u *UUID) ToString() string {
	return u.UUID.String()
}

// IsZero returns true if the UUID is equal to the zero-value.
func (u *UUID) IsZero() bool {
	return u.UUID == gid.UUID([gid.Size]byte{0})
}

// Equal returns true if both argument UUIDs contain the same value or if both are nil.
func Equal(a *UUID, b *UUID) bool {
	if a == nil || b == nil {
		return a == b
	}

	return gid.Equal(a.UUID, b.UUID)
}
