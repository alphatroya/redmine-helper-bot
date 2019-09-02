package mocks

import "fmt"

type StorageMock struct {
	storageToken map[int64]string
	storageHost  map[int64]string
}

func (r *StorageMock) ResetData(chat int64) error {
	r.storageHost[chat] = ""
	r.storageToken[chat] = ""
	return nil
}

func (r *StorageMock) StorageHost() map[int64]string {
	return r.storageHost
}

func (r StorageMock) StorageToken() map[int64]string {
	return r.storageToken
}

func NewStorageMock() *StorageMock {
	mock := new(StorageMock)
	mock.storageToken = make(map[int64]string)
	mock.storageHost = make(map[int64]string)
	return mock
}

func (r StorageMock) SetToken(token string, chat int64) {
	r.storageToken[chat] = token
}

func (r StorageMock) GetToken(chat int64) (string, error) {
	token, ok := r.storageToken[chat]
	if !ok {
		return "", fmt.Errorf("storage value is nil")
	}
	return token, nil
}

func (r StorageMock) SetHost(host string, chat int64) {
	r.storageHost[chat] = host
}

func (r StorageMock) GetHost(chat int64) (string, error) {
	host, ok := r.storageHost[chat]
	if !ok {
		return "", fmt.Errorf("storage value is nil")
	}
	return host, nil
}
