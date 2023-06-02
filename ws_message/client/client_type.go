package client

type ClientType string

const (
	// JsonClientType json client
	JsonClientType ClientType = "json"
	// bytesClientType bytes client
	bytesClientType ClientType = "bytes"
	// stringClientType string client
	stringClientType ClientType = "string"
)
