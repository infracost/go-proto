package quicksight

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type User struct {
	resource.Resource `tree:"-"`
	UserRole          value.Value[UserRole] `tree:"user_role"`
}

type UserRole uint32

const (
	UserRoleUnknown UserRole = iota
	UserRoleReader
	UserRoleAuthor
	UserRoleAdmin
	UserRoleReaderPro
	UserRoleAuthorPro
	UserRoleAdminPro
	UserRoleRestrictedReader
	UserRoleRestrictedAuthor
)
