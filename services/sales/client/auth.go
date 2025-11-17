package client

import (
	"microservice-challenge/package/client"
)

type AuthClient = client.AuthClient

func NewAuthClient(baseURL, serviceName, serviceSecret string) *AuthClient {
	return client.NewAuthClient(baseURL, serviceName, serviceSecret)
}
