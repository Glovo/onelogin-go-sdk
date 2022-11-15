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

type clientAppsWrite struct {
	ID        *int32  `json:"app_id,omitempty"`
	APIAuthID *int32  `json:"api_auth_id,omitempty"`
	Name      *string `json:"name,omitempty"`
	Scopes    []int32 `json:"scopes,omitempty"`
}

func (c *ClientApp) asWrite() *clientAppsWrite {
	return &clientAppsWrite{
		ID:        c.AppID,
		APIAuthID: c.APIAuthID,
		Name:      c.Name,
		Scopes:    c.ScopeIDs,
	}
}
