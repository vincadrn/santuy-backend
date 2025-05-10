package model

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
	Role string
}

// assoc table: groupuser: id_user, id_group, role
type UserGroupRole struct {
	User
	Group
	Role string
}
