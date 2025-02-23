package model

import (
	"context"
	"database/sql"
	"log"
)

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

type GroupGroupRole struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type GroupGroupRoles []GroupGroupRole

func (u *User) GetUser(db *sql.DB, ctx context.Context, userId string) {
	err := db.QueryRowContext(
		ctx,
		// `SELECT "ID_User", "ID_Group", "Nama", "Email", "Role" FROM "User" WHERE "ID_User" = $1;`,
		`SELECT "ID_User", "Nama", "Email", "Role" FROM "User" WHERE "ID_User" = $1;`,
		userId,
	).Scan()

	if err != nil {
		log.Fatal(err)
	}
}

func (u User) ToJSON() string {
	return MarshalToJSON(u)
}

func (g Group) ToJSON() string {
	return MarshalToJSON(g)
}

func (ggr GroupGroupRole) ToJSON() string {
	return MarshalToJSON(ggr)
}

func (ggrs GroupGroupRoles) ToJSON() string {
	return MarshalToJSON(ggrs)
}
