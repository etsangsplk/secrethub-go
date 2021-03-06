package api

// Error
var (
	ErrAccessLevelUnknown = errAPI.Code("access_level_unknown").Error("The access level is not known")
)

// Permission defines what kind of access an access rule grants or a access level has.
type Permission int

// The different Permission options.
const (
	PermissionNone Permission = iota
	PermissionRead
	PermissionWrite
	PermissionAdmin
)

// Set sets the Permission to the value.
func (al *Permission) Set(value string) error {
	switch value {
	case "r", "read":
		*al = PermissionRead
	case "w", "write":
		*al = PermissionWrite
	case "a", "admin":
		*al = PermissionAdmin
	case "n", "none":
		*al = PermissionNone
	default:
		return ErrAccessLevelUnknown
	}
	return nil
}

func (al Permission) String() string {
	switch al {
	case PermissionRead:
		return "read"
	case PermissionWrite:
		return "write"
	case PermissionAdmin:
		return "admin"
	default:
		return "none"
	}
}
