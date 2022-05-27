.PHONY: openapi
openapi:
	oapi-codegen -generate chi-server -o pkg/interfaces/openapi/openapi_api_gen.go -package openapi api/policy.yml
	oapi-codegen -generate types -o pkg/interfaces/openapi/openapi_types_gen.go -package openapi api/policy.yml