package response

import "vincadrn.com/santuy/internal/model"

type User struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type GroupRole struct {
	GroupName string `json:"group_name"`
	Role      string `json:"role"`
}

type GroupRoles []GroupRole

type UserGroupRole struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Groups   []GroupRole
}

func (res User) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}

func (res GroupRole) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}

func (res *GroupRole) Construct(groupRole *model.GroupRole) {
	res.GroupName = groupRole.Name
	res.Role = groupRole.Role
}

func (ress GroupRoles) ToJSON() ([]byte, error) {
	return MarshalToJSON(ress)
}

// Construct from model
func (ress *GroupRoles) Construct(groupRoles *[]model.GroupRole) {
	for _, val := range *groupRoles {
		var groupRole GroupRole
		groupRole.Construct(&val)

		*ress = append(*ress, groupRole)
	}
}

func (res UserGroupRole) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}
