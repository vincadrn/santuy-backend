package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"vincadrn.com/santuy/internal/model"
	"vincadrn.com/santuy/internal/response"
)

func TestGroupRoleResponse(t *testing.T) {
	role := model.GroupRole{
		Group: model.Group{
			Id:   "1",
			Name: "group1",
		},
		Role: model.Admin,
	}

	resp := response.GroupRole{}
	resp.Construct(&role)
	_, err := resp.ToJSON()

	assert.Equal(t, "group1", resp.GroupName, "group name should be equal")
	assert.Equal(t, "admin", resp.Role, "role should be correct")
	assert.NoError(t, err, "json marshalling should not produce error")
}

func TestGroupRolesResponse(t *testing.T) {
	roles := []model.GroupRole{
		{
			Group: model.Group{
				Id:   "1",
				Name: "group1",
			},
			Role: model.Admin,
		},
		{
			Group: model.Group{
				Id:   "2",
				Name: "group2",
			},
			Role: model.Member,
		},
	}

	resp := response.GroupRoles{}
	resp.Construct(&roles)
	_, err := resp.ToJSON()

	assert.Equal(t, "group1", resp[0].GroupName, "group name should be equal")
	assert.Equal(t, "admin", resp[0].Role, "role should be correct")
	assert.Equal(t, "group2", resp[1].GroupName, "group name should be equal")
	assert.Equal(t, "member", resp[1].Role, "role should be correct")
	assert.NoError(t, err, "json marshalling should not produce error")
}
