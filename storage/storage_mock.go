package storage

import "fmt"

type Mock struct {
	storageToken map[int64]string
	storageHost  map[int64]string
}

func (r *Mock) ResetData(chat int64) error {
	r.storageHost[chat] = ""
	r.storageToken[chat] = ""
	return nil
}

func (r *Mock) StorageHost() map[int64]string {
	return r.storageHost
}

func (r Mock) StorageToken() map[int64]string {
	return r.storageToken
}

func NewStorageMock() *Mock {
	mock := new(Mock)
	mock.storageToken = make(map[int64]string)
	mock.storageHost = make(map[int64]string)
	return mock
}

func (r Mock) SetToken(token string, chat int64) {
	r.storageToken[chat] = token
}

func (r Mock) GetToken(chat int64) (string, error) {
	token, ok := r.storageToken[chat]
	if !ok {
		return "", fmt.Errorf("storage value is nil")
	}
	return token, nil
}

func (r Mock) SetHost(host string, chat int64) {
	r.storageHost[chat] = host
}

func (r Mock) GetHost(chat int64) (string, error) {
	host, ok := r.storageHost[chat]
	if !ok {
		return "", fmt.Errorf("storage value is nil")
	}
	return host, nil
}
