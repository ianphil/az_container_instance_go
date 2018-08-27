package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

func main() {
	containerClient := containerinstance.NewContainerGroupsClient("2295f62b-34e7-40a1-9e9f-6def6b9f20b7")
	//containerClient.RequestInspector = logRequest()
	//containerClient.ResponseInspector = logResponse()

	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		containerClient.Authorizer = authorizer
		fmt.Println("Authorized")
	}

	containerProps := containerinstance.ContainerProperties{
		Image: to.StringPtr("nginx:latest"),
		Ports: &[]containerinstance.ContainerPort{
			containerinstance.ContainerPort{
				Protocol: containerinstance.ContainerNetworkProtocolTCP,
				Port:     to.Int32Ptr(80),
			},
		},
		Resources: &containerinstance.ResourceRequirements{
			Requests: &containerinstance.ResourceRequests{
				MemoryInGB: to.Float64Ptr(2),
				CPU:        to.Float64Ptr(1),
			},
		},
	}

	containerGroup := containerinstance.ContainerGroup{
		//Name:     to.StringPtr("mrjenk"),
		Location: to.StringPtr("eastus2"),
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			Containers: &[]containerinstance.Container{
				containerinstance.Container{
					Name:                to.StringPtr("jenkins"),
					ContainerProperties: &containerProps,
				},
			},
			RestartPolicy: containerinstance.Always,
			OsType:        containerinstance.Linux,
			IPAddress: &containerinstance.IPAddress{
				Ports: &[]containerinstance.Port{
					containerinstance.Port{
						Protocol: "TCP",
						Port:     to.Int32Ptr(80),
					},
				},
				Type:         to.StringPtr("Public"),
				DNSNameLabel: to.StringPtr("ipcontnginx"),
			},
		},
	}

	_, err = containerClient.CreateOrUpdate(context.Background(),
		"dxc_exp",
		"mrjenk",
		containerGroup)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Completed")
}

func logRequest() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpRequestOut(r, true)
			log.Println(string(dump))
			return r, err
		})
	}
}

func logResponse() autorest.RespondDecorator {
	return func(p autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(r *http.Response) error {
			err := p.Respond(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpResponse(r, true)
			log.Println(string(dump))
			return err
		})
	}
}
