package api

//go:generate oapi-codegen -o ../internal/controller/http/gen/types.go -generate types -package gen schema.yaml
//go:generate oapi-codegen  -o ../internal/controller/http/gen/server.go -generate chi-server,strict-server -package gen schema.yaml
//go:generate oapi-codegen  -o ../internal/controller/http/gen/spec.go -generate spec -package gen schema.yaml
