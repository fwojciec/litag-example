package mocks

//go:generate moq -out q.go -pkg mocks ../../postgres Q
//go:generate moq -out txq.go -pkg mocks ../../postgres TxQ
