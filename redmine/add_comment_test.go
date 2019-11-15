package redmine

import (
	"testing"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

func TestClientManager_AddComment(t *testing.T) {
	for _, isError := range []bool{true, false} {
		var networkClient *ClientRequestMock
		if isError {
			networkClient = NewClientRequestMock(400, nil, "")
		} else {
			networkClient = NewClientRequestMock(200, nil, "")
		}
		storageMock := storage.NewStorageMock()
		storageMock.SetHost("http://google.com", 10)
		storageMock.SetToken("fdjsdfjs", 10)
		redmine := NewClientManager(networkClient, storageMock, 10)
		err := redmine.AddComment("4333", "FooBar")
		if isError && err == nil {
			t.Errorf("add comment not return error but it should")
		} else if !isError && err != nil {
			t.Errorf("add comment return error: %s, but it shouldn't", err)
		}
	}
}
