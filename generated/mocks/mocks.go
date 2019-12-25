package mocks

//go:generate moq -out querent.go -pkg mocks ../../postgres Querent
//go:generate moq -out txquerent.go -pkg mocks ../../postgres TxQuerent
