package discovery

import (
	"testing"
)

func TestParseMetaData(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{"good json", "{\n    \"authorization_endpoint\": \"https://wallet.hello.coop/authorize\",\n    \"issuer\": \"https://issuer.hello.coop\",\n    \"jwks_uri\": \"https://issuer.hello.coop/.well-known/jwks\",\n    \"introspection_endpoint\": \"https://wallet.hello.coop/oauth/introspect\",\n    \"token_endpoint\": \"https://wallet.hello.coop/oauth/token\",\n    \"userinfo_endpoint\": \"https://wallet.hello.coop/oauth/userinfo\",\n    \"service_documentation\": \"https://hello.dev\",\n    \"response_modes_supported\": [\n        \"query\",\n        \"fragment\",\n        \"form_post\"\n        ],\n    \"subject_types_supported\": [\n        \"pairwise\"\n        ],\n    \"id_token_signing_alg_values_supported\": [\n        \"RS256\"\n        ],\n    \"token_endpoint_auth_methods_supported\": [\n        \"client_secret_basic\"\n        ],\n    \"introspection_endpoint_auth_methods_supported\": [\n        \"none\"\n        ],\n    \"code_challenge_methods_supported\": [\n        \"S256\"\n        ],\n    \"grant_types_supported\": [\n        \"authorization_code\", \n        \"implicit\"\n        ],\n    \"response_types_supported\": [\n        \"id_token\",\n        \"code\"\n        ],\n    \"scopes_supported\": [\n        \"openid\",\n        \"name\",\n        \"nickname\",\n        \"family_name\",\n        \"given_name\",\n        \"picture\",\n        \"email\",\n        \"phone\",\n        \"profile_update\",\n        \"ethereum\"\n        ],\n    \"claims_supported\": [\n        \"sub\",\n        \"iss\",\n        \"aud\",\n        \"exp\",\n        \"iat\",\n        \"jti\",\n        \"nonce\",\n        \"name\",\n        \"picture\",\n        \"email\",\n        \"email_verified\",\n        \"phone\",\n        \"phone_verified\",\n        \"ethereum\"\n        ]\n}", false},
		{"bad json", "{\n    \"authorization_endpoint\": \"https://wallet.hello.coop/authorize\",\n    \"issuer\": \"https://issuer.hello.coop\",\n    \"jwks_uri\": \"https://issuer.hello.coop/.well-known/jwks\",\n    \"introspection_endpoint\": \"https://wallet.hello.coop/oauth/introspect\",\n    \"token_endpoint\": \"https://wallet.hello.coop/oauth/token\",\n    \"userinfo_endpoint\": \"https://wallet.hello.coop/oauth/userinfo\",\n    \"service_documentation\": \"https://hello.dev\",\n    \"response_modes_supported\": [\n        \"query\",\n        \"fragment\",\n        \"form_post\"\n        ],\n    \"subject_types_supported\": [\n        \"pairwise\"\n        ],\n    \"id_token_signing_alg_values_supported\": [\n        \"RS256\"\n        ],\n    \"token_endpoint_auth_methods_supported\": [\n        \"client_secret_basic\"\n        ],\n    \"introspection_endpoint_auth_methods_supported\": [\n        \"none\"\n        ],\n    \"code_challenge_methods_supported\": [\n        \"S256\"\n        ],\n    \"grant_types_supported\": [\n        \"authorization_code\", \n        \"implicit\"\n        ],\n    \"response_types_supported\": [\n        \"id_token\",\n        \"code\"\n        ],\n    \"scopes_supported\": [\n        \"openid\",\n        \"name\",\n        \"nickname\",\n        \"family_name\",\n        \"given_name\",\n        \"picture\",\n        \"email\",\n        \"phone\",\n        \"profile_update\",\n        \"ethereum\"\n        ],\n    \"claims_supported\": [\n        \"sub\",\n        \"iss\",\n        \"aud\",\n        \"exp\",\n        \"iat\",\n        \"jti\",\n        \"nonce\",\n        \"name\",\n        \"picture\",\n        \"email\",\n        \"email_verified\",\n        \"phone\",\n        \"phone_verified\",\n        \"ethereum\"\n        ", true},
		{"no json", "", true},
		{"empty", "{}", false},
		// TODO in config allow to change size of unreasonable length
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseMetaData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMetaData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
