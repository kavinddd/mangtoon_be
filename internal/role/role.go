package role

// constant values of roles in DB
const (
	Reader = "reader"
	Writer = "writer"
	Admin  = "admin"
)

func IsAdmin(roles []string) bool {
	return IsRole(roles, Admin)
}

func IsWriter(roles []string) bool {
	return IsRole(roles, Writer)
}

func IsReader(roles []string) bool {
	return IsRole(roles, Reader)
}

func IsRole(roles []string, theRole string) bool {
	for _, role := range roles {
		if role == theRole {
			return true
		}
	}

	return false
}
