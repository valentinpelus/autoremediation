// En cours de rédaction
package main

import (
	"regexp"
	"testing"

	client "github.com/kubernetes-sdk-for-go-101/pkg/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

type Client struct {
	Clientset kubernetes.Interface
}

// MockClient is the mock client
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do is the mock client's `Do` func
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func Test_getVMAlertBackendSize(t *testing.T) {
	var client client.Client
	client.Clientset = testclient.NewSimpleClientset()
	type args struct {
		server string
		Clientset *KUBERNETES.Clientset
	}
	tests := []struct {
		name string
		args args
		want map[string]bool
		response func(*http.Request) (*http.Response, error)
	}{
		name: "Alert is firing",
		args: args{
			server: "https://some.server",
			clientset clientset
		}
	}
}

