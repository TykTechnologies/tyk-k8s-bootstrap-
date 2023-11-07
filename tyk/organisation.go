package tyk

import (
	"bytes"
	"errors"
	"io"

	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	ic "tyk/tyk/bootstrap/pkg/constants"
	"tyk/tyk/bootstrap/tyk/api"
	"tyk/tyk/bootstrap/tyk/internal/constants"
)

var ErrOrgExists = errors.New("there shouldn't be any organisations, please " +
	"disable bootstrapping to avoid losing data or delete " +
	"already existing organisations")

// OrgExists checks if the given Tyk Organisation is created or not.
// It returns ErrOrgExists if the organisation exists.
func (s *Service) OrgExists() error {
	orgsApiEndpoint := s.appArgs.K8s.DashboardSvcUrl + constants.AdminOrganisationsEndpoint

	req, err := http.NewRequest(http.MethodGet, orgsApiEndpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set(ic.AdminAuthHeader, s.appArgs.Tyk.Admin.Secret)
	req.Header.Set(ic.ContentTypeHeader, "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	orgs := api.OrgAPIResp{}

	err = json.Unmarshal(bodyBytes, &orgs)
	if err != nil {
		return err
	}

	if len(orgs.Organisations) > 0 {
		for _, organisation := range orgs.Organisations {
			if organisation["owner_name"] == s.appArgs.Tyk.Org.Name ||
				organisation["cname"] == s.appArgs.Tyk.Org.Cname {
				return ErrOrgExists
			}
		}
	}

	return nil
}

// CreateOrganisation creates organisation based on the information populated in the config.Config.
func (s *Service) CreateOrganisation() error {
	createOrgData := api.CreateOrgReq{
		OwnerName:    s.appArgs.Tyk.Org.Name,
		CnameEnabled: true,
		Cname:        s.appArgs.Tyk.Org.Cname,
	}

	reqBodyBytes, err := json.Marshal(createOrgData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.appArgs.K8s.DashboardSvcUrl+constants.AdminOrganisationsEndpoint,
		bytes.NewReader(reqBodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set(ic.AdminAuthHeader, s.appArgs.Tyk.Admin.Secret)
	req.Header.Set(ic.ContentTypeHeader, "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	createOrgResp := api.DashboardAPIResp{}

	err = json.Unmarshal(bodyBytes, &createOrgResp)
	if err != nil {
		return err
	}

	s.appArgs.Tyk.Org.ID = createOrgResp.Meta

	return nil
}
