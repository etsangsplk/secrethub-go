package secrethub

import (
	"github.com/secrethub/secrethub-go/internals/api"
	"github.com/secrethub/secrethub-go/internals/errio"
	"github.com/secrethub/secrethub-go/pkg/secrethub/credentials"
)

// ServiceService handles operations on service accounts from SecretHub.
type ServiceService interface {
	// Create creates a new service account for the given repo.
	Create(path string, description string, credential credentials.Creator) (*api.Service, error)
	// Delete removes a service account by name.
	Delete(name string) (*api.RevokeRepoResponse, error)
	// Get retrieves a service account by name.
	Get(name string) (*api.Service, error)
	// List lists all service accounts in a given repository.
	List(path string) ([]*api.Service, error)
}

func newServiceService(client *Client) ServiceService {
	return serviceService{
		client: client,
	}
}

type serviceService struct {
	client *Client
}

// Create creates a new service account for the given repo.
func (s serviceService) Create(path string, description string, credentialCreator credentials.Creator) (*api.Service, error) {
	repoPath, err := api.NewRepoPath(path)
	if err != nil {
		return nil, errio.Error(err)
	}

	err = api.ValidateServiceDescription(description)
	if err != nil {
		return nil, errio.Error(err)
	}

	verifier, encrypter, metadata, err := credentialCreator.Create()
	if err != nil {
		return nil, err
	}

	accountKey, err := generateAccountKey()
	if err != nil {
		return nil, errio.Error(err)
	}

	credentialRequest, err := s.client.createCredentialRequest(verifier, metadata)
	if err != nil {
		return nil, errio.Error(err)
	}

	accountKeyRequest, err := s.client.createAccountKeyRequest(encrypter, accountKey)
	if err != nil {
		return nil, errio.Error(err)
	}

	serviceRepoMemberRequest, err := s.client.createRepoMemberRequest(repoPath, accountKeyRequest.PublicKey)
	if err != nil {
		return nil, errio.Error(err)
	}

	in := &api.CreateServiceRequest{
		Description: description,
		Credential:  credentialRequest,
		AccountKey:  accountKeyRequest,
		RepoMember:  serviceRepoMemberRequest,
	}

	err = in.Validate()
	if err != nil {
		return nil, errio.Error(err)
	}

	service, err := s.client.httpClient.CreateService(repoPath.GetNamespace(), repoPath.GetRepo(), in)
	if err != nil {
		return nil, errio.Error(err)
	}

	return service, nil
}

// Delete removes a service account.
func (s serviceService) Delete(name string) (*api.RevokeRepoResponse, error) {
	err := api.ValidateServiceID(name)
	if err != nil {
		return nil, errio.Error(err)
	}

	resp, err := s.client.httpClient.DeleteService(name)
	if err != nil {
		return nil, errio.Error(err)
	}

	return resp, nil
}

// Get retrieves a service account.
func (s serviceService) Get(name string) (*api.Service, error) {
	err := api.ValidateServiceID(name)
	if err != nil {
		return nil, errio.Error(err)
	}

	return s.client.httpClient.GetService(name)
}

// List is an alias of the RepoServiceService List function.
func (s serviceService) List(path string) ([]*api.Service, error) {
	repoServiceService := newRepoServiceService(s.client)
	return repoServiceService.List(path)
}
