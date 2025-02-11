package api

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o ../internal/controller/http/gen/types.go -generate types -package gen schema.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o ../internal/controller/http/gen/server.go -generate chi-server,strict-server -package gen schema.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o ../internal/controller/http/gen/spec.go -generate spec -package gen schema.yaml
