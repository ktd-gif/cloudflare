package cloudflare

import (
	"strings"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/cloudflare"
)

var benchTokens = map[string]string{
	"legacy":       "Sqqty8-Vn0iOP29rvqYgwKz_xqGQ4y5JhuVL1-qU",
	"user":         "cfut_6YNSv6zUiWehDxZfdh8vNf3orQJzA6rrcDFJhql65e203820",
	"account":      "cfat_3Sak5MsJJAsXNHCEmRufjbkbrdcxrhZFofmAe1pN435b71a9",
	"invalid":      "{env.CF_API_TOKEN}",
	"too_long_new": "cfut_" + strings.Repeat("a", 257),
}

var benchSink bool

func BenchmarkValidCloudflareToken(b *testing.B) {
	for _, name := range []string{"legacy", "user", "account", "invalid", "too_long_new"} {
		token := benchTokens[name]
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchSink = validCloudflareToken(token)
			}
		})
	}
}

func BenchmarkProvision(b *testing.B) {
	for _, name := range []string{"legacy", "user"} {
		token := benchTokens[name]
		b.Run(name, func(b *testing.B) {
			p := Provider{&cloudflare.Provider{}}
			for i := 0; i < b.N; i++ {
				p.Provider.APIToken = token
				if err := p.Provision(caddy.Context{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkUnmarshalCaddyfile(b *testing.B) {
	configs := []struct {
		name   string
		config string
	}{
		{"single_arg", "cloudflare Sqqty8-Vn0iOP29rvqYgwKz_xqGQ4y5JhuVL1-qU"},
		{"api_token_block", `cloudflare {
			api_token Sqqty8-Vn0iOP29rvqYgwKz_xqGQ4y5JhuVL1-qU
		}`},
		{"zone_and_api_tokens", `cloudflare {
			zone_token cfut_6YNSv6zUiWehDxZfdh8vNf3orQJzA6rrcDFJhql65e203820
			api_token cfat_3Sak5MsJJAsXNHCEmRufjbkbrdcxrhZFofmAe1pN435b71a9
		}`},
	}

	for _, tc := range configs {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d := caddyfile.NewTestDispenser(tc.config)
				p := Provider{&cloudflare.Provider{}}
				if err := p.UnmarshalCaddyfile(d); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
