package model

type User struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Email  string      `json:"email"`
	Groups []GroupRole `json:"groups"`
}

type Group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GroupRole struct {
	GroupId string `json:"group_id"`
	Role    string `json:"role"`
}

func (u User) ToJSON() string {
	return MarshalToJSON(u)
}

func (g Group) ToJSON() string {
	return MarshalToJSON(g)
}
