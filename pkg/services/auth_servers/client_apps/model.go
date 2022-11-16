package clientapps

import "github.com/glovo/onelogin-go-sdk/pkg/services/auth_servers/scopes"

type ClientAppsQuery struct {
	AuthServerID string
}

type ClientApp struct {
	AppID     *int32         `json:"app_id,omitempty"`
	APIAuthID *int32         `json:"api_auth_id,omitempty"`
	Name      *string        `json:"name,omitempty"`
	Scopes    []scopes.Scope `json:"scopes,omitempty"`
	ScopeIDs  []int32        `json:"scope_ids,omitempty"`
}

type clientAppsCreate struct {
	ID        *int32  `json:"app_id,omitempty"`
	APIAuthID *int32  `json:"api_auth_id,omitempty"`
	Scopes    []int32 `json:"scopes,omitempty"`
}

type clientAppsUpdate struct {
	Scopes []int32 `json:"scopes,omitempty"`
}

func (c *ClientApp) asCreate() *clientAppsCreate {
	return &clientAppsCreate{
		ID:        c.AppID,
		APIAuthID: c.APIAuthID,
		Scopes:    c.ScopeIDs,
	}
}

func (c *ClientApp) asUpdate() *clientAppsUpdate {
	return &clientAppsUpdate{
		Scopes: c.ScopeIDs,
	}
}
