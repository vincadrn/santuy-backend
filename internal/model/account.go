package model

type RoleType int

const (
	Admin RoleType = iota
	Member
)

var roleType = map[RoleType]string{
	Admin:  "admin",
	Member: "member",
}

// RoleType implements Stringer
func (rt RoleType) String() string {
	return roleType[rt]
}

// user_table
type User struct {
	Id    string // id_user
	Name  string // nama
	Email string // email
}

// group_table
type Group struct {
	Id   string // id_group
	Name string // nama
}

type GroupRole struct {
	Group
	Role RoleType
}

// assoc table: groupuser: id_user, id_group, role
type UserGroupRole struct {
	User
	Group
	Role string
}
