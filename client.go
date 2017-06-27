// package rpc represent the remote procedure call 
//which import the package rpc and create a new api client for docker-volume-client 
package rpc

import "net/http"

type rpcClient struct {
	client   *http.Client
	url      string
	username string
	password string
}

// NewrpcClient creates a new docker-volume-dirverv API client
func NewrpcClient(url string, username string, password string) *rpcClient {
	return &rpcClient{
		client:   &http.Client{},
		url:      url,
		username: username,
		password: password,
	}
}

// CreateVolume creates a new volume. Its root directory will be owned by given user and group
func (client rpcClient) CreateVolume(request *CreateVolumeRequest) (string, error) {
	var response volumeUUID
	if err := client.sendRequest("createVolume", request, &response); err != nil {
		return "", err
	}

	return response.VolumeUUID, nil
	//returns a uuid of the created volume 
}

// ResolveVolumeNameToUUID resolves a volume name to a UUID
func (client *rpcClient) ResolveVolumeNameToUUID(volumeName, tenant string) (string, error) {
	request := &resolveVolumeNameRequest{
		VolumeName:   volumeName,
		TenantDomain: tenant,
	}
	var response volumeUUID
	if err := client.sendRequest("resolveVolumeName", request, &response); err != nil {
		return "", err
	}

	return response.VolumeUUID, nil
}

// DeleteVolume deletes a docker-volume-driver volume
func (client *rpcClient) DeleteVolume(UUID string) error {
	return client.sendRequest(
		"deleteVolume",
		&volumeUUID{
			VolumeUUID: UUID,
		},
		nil)
}

// DeleteVolumeByName deletes a volume by a given name
func (client *rpcClient) DeleteVolumeByName(volumeName, tenant string) error {
	uuid, err := client.ResolveVolumeNameToUUID(volumeName, tenant)
	if err != nil {
		return err
	}

	return client.DeleteVolume(uuid)
}

// GetClientList returns a list of all active clients
func (client *rpcClient) GetClientList(tenant string) (GetClientListResponse, error) {
	request := &getClientListRequest{
		TenantDomain: tenant,
	}

	var response GetClientListResponse
	if err := client.sendRequest("getClientListRequest", request, &response); err != nil {
		return response, err
	}

	return response, nil
}
