package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// GeneratePolicyExamplesParams generate_policy_examples工具的参数
type GeneratePolicyExamplesParams struct {
	ConfigType string `json:"config_type"` // http_forward|sni_forward|socks_forward|web_doh|dns
	PolicyType string `json:"policy_type"` // geoip|geosite|domain_match|ip_range|file_based|fetch_based|custom
}

// GeneratePolicyExamples 生成Policy模板示例和文档
func GeneratePolicyExamples(arguments json.RawMessage) (string, error) {
	var params GeneratePolicyExamplesParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid config_type and policy_type",
		)
	}

	log.Info().
		Str("config_type", params.ConfigType).
		Str("policy_type", params.PolicyType).
		Msg("generating policy examples")

	// 设置默认值
	if params.ConfigType == "" {
		params.ConfigType = "http_forward"
	}
	if params.PolicyType == "" {
		params.PolicyType = "geoip"
	}

	var content string

	// 生成对应的Policy模板示例
	switch params.PolicyType {
	case "geoip":
		content = generateGeoIPPolicyExample(params.ConfigType)
	case "geosite":
		content = generateGeositePolicyExample(params.ConfigType)
	case "domain_match":
		content = generateDomainMatchPolicyExample(params.ConfigType)
	case "ip_range":
		content = generateIPRangePolicyExample(params.ConfigType)
	case "file_based":
		content = generateFileBasedPolicyExample(params.ConfigType)
	case "fetch_based":
		content = generateFetchBasedPolicyExample(params.ConfigType)
	case "custom":
		content = generateCustomPolicyExample(params.ConfigType)
	default:
		return responses.ErrorResponse(
			fmt.Sprintf("Unknown policy_type: %s", params.PolicyType),
			"Supported types: geoip, geosite, domain_match, ip_range, file_based, fetch_based, custom",
		)
	}

	description := fmt.Sprintf("Policy template example for %s with %s routing", params.ConfigType, params.PolicyType)
	log.Info().Msg("policy examples generated successfully")
	return responses.SuccessResponse(content, description)
}

func generateGeoIPPolicyExample(configType string) string {
	return `# GeoIP-based Routing Policy
# Routes traffic based on the geographic location of the destination IP

# Example 1: Route China IP to direct, others to proxy
{{ if eq (country (dnsResolve .Request.Host)) "CN" }}direct{{ else }}proxy{{ end }}

# Example 2: Multiple country routing
{{ $country := country (dnsResolve .Request.Host) }}
{{ if eq $country "CN" }}cn_dialer
{{ else if eq $country "US" }}us_dialer
{{ else }}default_dialer{{ end }}

# Example 3: Using geoip function for detailed info
{{ $geoip := geoip (dnsResolve .Request.Host) }}
{{ if eq $geoip.Country "CN" }}direct
{{ else if eq $geoip.ISP "Amazon" }}aws_dialer
{{ else }}proxy{{ end }}

# Available GeoIP fields:
# - .Country: Country code (e.g., "CN", "US")
# - .City: City name
# - .ISP: Internet Service Provider
# - .ConnectionType: Connection type

# Available template functions:
# - country(ip): Get country code from IP
# - geoip(ip): Get full GeoIP information
# - dnsResolve(host): Resolve hostname to IP
`
}

func generateGeositePolicyExample(configType string) string {
	return `# Geosite-based Routing Policy
# Routes traffic based on domain categorization (similar to V2Ray routing)

# Example 1: Route Chinese sites directly
{{ if eq (geosite .Request.Host) "cn" }}direct{{ else }}proxy{{ end }}

# Example 2: Multiple site categories
{{ $site := geosite .Request.Host }}
{{ if eq $site "cn" }}direct
{{ else if eq $site "google" }}google_dialer
{{ else if eq $site "netflix" }}netflix_dialer
{{ else }}proxy{{ end }}

# Example 3: Combined with domain
{{ $domain := domain .Request.Host }}
{{ $site := geosite .Request.Host }}
{{ if eq $site "cn" }}direct
{{ else if hasSuffixes "google.com|youtube.com" $domain }}google_dialer
{{ else }}proxy{{ end }}

# Common geosite categories:
# - cn: Chinese websites
# - google: Google services
# - netflix: Netflix services
# - category-ads: Advertisement domains
# - geolocation-!cn: Non-Chinese sites

# Available template functions:
# - geosite(domain): Get domain category
# - domain(host): Extract root domain
# - host(hostport): Extract hostname from host:port
`
}

func generateDomainMatchPolicyExample(configType string) string {
	return `# Domain Matching Policy
# Routes traffic based on domain name patterns

# Example 1: Suffix matching
{{ if hasSuffixes "google.com|youtube.com|googleapis.com" .Request.Host }}google_dialer{{ else }}proxy{{ end }}

# Example 2: Prefix matching
{{ if hasPrefixes "api.|cdn." .Request.Host }}api_dialer{{ else }}direct{{ end }}

# Example 3: Wildcard matching
{{ if wildcardMatch "*.google.com|*.youtube.com" .Request.Host }}google_dialer{{ else }}proxy{{ end }}

# Example 4: Regex matching
{{ if regexMatch "^(www|api)\\.(google|youtube)\\.com$" .Request.Host }}google_dialer{{ else }}proxy{{ end }}

# Example 5: Extract and match root domain
{{ $domain := domain .Request.Host }}
{{ if eq $domain "google.com" }}google_dialer
{{ else if eq $domain "facebook.com" }}facebook_dialer
{{ else }}proxy{{ end }}

# Example 6: Multiple conditions
{{ $host := host .Request.Host }}
{{ if hasSuffixes "google.com|youtube.com" $host }}google_dialer
{{ else if hasSuffixes "cn|com.cn" $host }}direct
{{ else }}proxy{{ end }}

# Available template functions:
# - hasSuffixes(pattern, str): Check if str has any suffix in pattern (pipe-separated)
# - hasPrefixes(pattern, str): Check if str has any prefix in pattern
# - wildcardMatch(pattern, str): Wildcard matching (* and ?)
# - regexMatch(pattern, str): Regular expression matching
# - domain(host): Extract root domain (e.g., "www.google.com" -> "google.com")
# - host(hostport): Extract hostname from "host:port"
`
}

func generateIPRangePolicyExample(configType string) string {
	return `# IP Range Routing Policy
# Routes traffic based on IP address ranges

# Example 1: Check if IP is in private network
{{ $ip := dnsResolve .Request.Host }}
{{ if isInNet $ip "10.0.0.0/8" }}direct
{{ else if isInNet $ip "192.168.0.0/16" }}direct
{{ else if isInNet $ip "172.16.0.0/12" }}direct
{{ else }}proxy{{ end }}

# Example 2: Route specific IP ranges to different dialers
{{ $ip := dnsResolve .Request.Host }}
{{ if isInNet $ip "1.0.0.0/8" }}apnic_dialer
{{ else if isInNet $ip "8.8.8.0/24" }}google_dialer
{{ else }}proxy{{ end }}

# Example 3: Using ipRange function for CIDR operations
{{ $range := ipRange "192.168.1.0/24" }}
{{ if $range.Contains .Request.Host }}direct{{ else }}proxy{{ end }}

# Example 4: Integer IP comparison
{{ $ipint := ipInt (dnsResolve .Request.Host) }}
{{ if and (ge $ipint 16777216) (le $ipint 33554431) }}direct{{ else }}proxy{{ end }}

# Available template functions:
# - isInNet(host, cidr): Check if host/IP is in CIDR range
# - dnsResolve(host): Resolve hostname to IP
# - ipRange(cidr): Get IP range object
# - ipInt(ip): Convert IPv4 to integer
`
}

func generateFileBasedPolicyExample(configType string) string {
	return `# File-based Routing Policy
# Routes traffic based on domain/IP lists stored in files

# Example 1: Check domain in file (one domain per line)
{{ if inFileLine "blocked_domains.txt" .Request.Host }}deny
{{ else if inFileLine "direct_domains.txt" .Request.Host }}direct
{{ else }}proxy{{ end }}

# Example 2: Check IP in file with IP sets/ranges
{{ $ip := dnsResolve .Request.Host }}
{{ if inFileIPSet "cn_ip.txt" $ip }}direct
{{ else if inFileIPSet "blocked_ip.txt" $ip }}deny
{{ else }}proxy{{ end }}

# Example 3: Combined domain and IP file checking
{{ $host := .Request.Host }}
{{ $ip := dnsResolve $host }}
{{ if inFileLine "whitelist_domains.txt" $host }}direct
{{ else if inFileIPSet "whitelist_ip.txt" $ip }}direct
{{ else if inFileLine "blacklist_domains.txt" $host }}deny
{{ else }}proxy{{ end }}

# File format for inFileLine (domains.txt):
# google.com
# youtube.com
# github.com

# File format for inFileIPSet (ips.txt):
# Supports: Single IPs, CIDR ranges, IP ranges
# 8.8.8.8
# 1.1.1.0/24
# 192.168.1.1-192.168.1.255

# Available template functions:
# - inFileLine(filename, line): Check if line exists in file (sorted search)
# - inFileIPSet(filename, ip): Check if IP is in file's IP set
# - dnsResolve(host): Resolve hostname to IP

# Note: Files are cached and auto-reloaded every 2 minutes
`
}

func generateFetchBasedPolicyExample(configType string) string {
	return `# Fetch-based Routing Policy
# Routes traffic based on dynamic data fetched from URLs

# Example 1: Fetch domain list from URL
{{ $resp := fetch "" 30 300 "https://example.com/blocked_domains.txt" }}
{{ if and (eq $resp.Status 200) (contains .Request.Host $resp.Body) }}deny{{ else }}proxy{{ end }}

# Example 2: Check domain in fetched list (line by line)
{{ $resp := fetch "" 30 300 "https://example.com/cn_domains.txt" }}
{{ if eq $resp.Status 200 }}
  {{ range $resp.Lines }}
    {{ if eq . $.Request.Host }}direct{{ end }}
  {{ end }}
{{ end }}
proxy

# Example 3: Fetch JSON and route based on data
{{ $resp := fetch "" 30 300 "https://api.example.com/routing.json" }}
{{ if eq $resp.Status 200 }}
  {{ $data := fromJson $resp.Body }}
  {{ $country := country (dnsResolve .Request.Host) }}
  {{ index $data $country }}
{{ else }}
  proxy
{{ end }}

# Example 4: Use custom User-Agent
{{ $resp := fetch "Mozilla/5.0" 30 300 "https://example.com/rules.txt" }}
{{ if contains .Request.Host $resp.Body }}special_dialer{{ else }}proxy{{ end }}

# fetch function signature:
# fetch(userAgent, timeoutSeconds, ttlSeconds, url) -> FetchResponse

# FetchResponse fields:
# - .Status: HTTP status code (int)
# - .Headers: HTTP response headers (map)
# - .Body: Response body as string
# - .Lines: Response body split by lines (for text/*)
# - .Error: Error if request failed
# - .CreatedAt: Timestamp of the fetch

# Available template functions:
# - fetch(ua, timeout, ttl, url): Fetch URL with caching
# - contains(substr, str): Check if str contains substr

# Note: Responses are cached for TTL seconds
`
}

func generateCustomPolicyExample(configType string) string {
	contextInfo := ""
	switch configType {
	case "http_forward", "socks_forward":
		contextInfo = `# Available context in templates:
# - .Request: *http.Request object
#   - .Request.Host: Target hostname
#   - .Request.Method: HTTP method
#   - .Request.URL: Request URL
#   - .Request.Header: HTTP headers
# - .RealIP: Client's real IP address (netip.Addr)
# - .ClientHelloInfo: TLS ClientHello information
# - .JA4: JA4 TLS fingerprint
# - .User: Authentication user info
#   - .User.Username: Username
#   - .User.Attrs: User attributes map
# - .UserAgent: Parsed user agent
#   - .UserAgent.OS: Operating system
#   - .UserAgent.Name: Browser name
#   - .UserAgent.Version: Browser version
# - .ServerAddr: Server listening address`
	case "sni_forward":
		contextInfo = `# Available context in templates:
# - .ServerName: SNI server name
# - .ClientHello: TLS ClientHello info
# - .RemoteAddr: Client address`
	case "web_doh":
		contextInfo = `# Available context in templates:
# - .Request: *http.Request object
# - .Question: DNS question
#   - .Question.Name: Query domain name
#   - .Question.Type: Query type (A, AAAA, etc.)`
	case "dns":
		contextInfo = `# Available context in templates:
# - .Question: DNS question
#   - .Question.Name: Query domain name
#   - .Question.Type: Query type
# - .RemoteAddr: Client address`
	}

	return fmt.Sprintf(`# Custom Policy Template Examples
# Advanced routing with multiple conditions and custom logic

# Example 1: User-based routing
{{ if eq .User.Username "admin" }}admin_dialer
{{ else if .User.Attrs.vip }}vip_dialer
{{ else }}normal_dialer{{ end }}

# Example 2: Time-based routing (using sprig functions)
{{ $hour := now | date "15" | atoi }}
{{ if and (ge $hour 9) (lt $hour 18) }}business_hours_dialer{{ else }}after_hours_dialer{{ end }}

# Example 3: Client IP-based routing
{{ $country := country .RealIP }}
{{ if eq $country "CN" }}cn_client_dialer
{{ else if eq $country "US" }}us_client_dialer
{{ else }}international_dialer{{ end }}

# Example 4: TLS fingerprint-based routing
{{ if .ClientHelloInfo }}
  {{ if greased .ClientHelloInfo }}browser_dialer{{ else }}bot_dialer{{ end }}
{{ else }}
  proxy
{{ end }}

# Example 5: Combined conditions
{{ $domain := domain .Request.Host }}
{{ $country := country (dnsResolve .Request.Host) }}
{{ if and (eq $domain "google.com") (eq $country "US") }}us_google_dialer
{{ else if eq $domain "google.com" }}google_dialer
{{ else if eq $country "CN" }}direct
{{ else }}proxy{{ end }}

# Example 6: UserAgent-based routing
{{ if eq .UserAgent.OS "iOS" }}mobile_dialer
{{ else if eq .UserAgent.Name "Chrome" }}chrome_dialer
{{ else }}default_dialer{{ end }}

# Example 7: IPv6 availability check
{{ if hasIPv6 .Request.Host }}ipv6_dialer{{ else }}ipv4_dialer{{ end }}

# Example 8: Special actions
# Return specific values for special behaviors:
# - "proxy_pass": Forward through proxy (default)
# - "direct": Direct connection
# - "deny" or "reject": Reject the request
# - "reset" or "close": Close connection immediately
# - "generate_204": Return HTTP 204 No Content
# - "require_auth": Require authentication
# - "bypass_auth": Bypass authentication check
# - "<dialer_name>": Use specific dialer

# Example 9: JSON-based dialer selection with options
{{ $country := country (dnsResolve .Request.Host) }}
{{ if eq $country "CN" }}{"dialer":"direct","disable_ipv6":true}
{{ else }}{"dialer":"proxy","prefer_ipv6":true}{{ end }}

# Example 10: Query string format
{{ if eq .User.Attrs.type "premium" }}dialer=premium&disable_ipv6=false{{ else }}dialer=normal{{ end }}

%s

# All Liner Template Functions:

## Go Template Built-in Functions:
# Control Structures:
# - if/else/end: Conditional execution
#   {{ if .condition }}...{{ else }}...{{ end }}
# - range/end: Iterate over arrays, slices, maps
#   {{ range .items }}{{ . }}{{ end }}
# - with/end: Set dot context
#   {{ with .field }}{{ . }}{{ end }}
# - define/template: Define and execute templates
#   {{ define "name" }}...{{ end }}
#   {{ template "name" . }}
# - block: Define and execute block
#   {{ block "name" . }}default{{ end }}
#
# Comparison & Logic:
# - eq, ne, lt, le, gt, ge: Comparison operators
#   {{ if eq .x .y }}equal{{ end }}
# - and, or, not: Logical operators
#   {{ if and .x .y }}both true{{ end }}
#
# Functions:
# - call: Call a function
#   {{ call .func .arg1 .arg2 }}
# - index: Index into maps, slices, arrays
#   {{ index .map "key" }}
# - len: Length of strings, slices, maps
#   {{ len .items }}
# - print, printf, println: Formatting output
#   {{ printf "%s: %d" .name .count }}

## Sprig Library Functions (from slim-sprig):
#
# String Functions:
# - trim, trimAll, trimPrefix, trimSuffix
# - upper, lower, title, untitle, snakecase, camelcase
# - repeat, substr, nospace, indent, quote, squote
# - contains, hasPrefix, hasSuffix, cat, replace
# - split, splitList, join, toString
#
# String List Functions:
# - sortAlpha, join, split
#
# Integer Math:
# - add, sub, mul, div, mod, max, min
# - ceil, floor, round
#
# Integer Array:
# - until, untilStep
#   {{ range until 5 }}{{ . }}{{ end }}  // 0 1 2 3 4
#
# Float Math:
# - addf, subf, mulf, divf, maxf, minf
#
# Date Functions:
# - now: Current time
#   {{ now | date "2006-01-02" }}
# - date: Format time
#   {{ .time | date "15:04:05" }}
# - dateModify: Modify date
#   {{ now | dateModify "+1h" }}
# - duration: Parse duration string
# - durationRound: Round duration
#
# Defaults:
# - default: Default value if empty
#   {{ .value | default "default" }}
# - empty: Check if empty
#   {{ if empty .value }}empty{{ end }}
# - coalesce: First non-empty value
#   {{ coalesce .a .b .c "default" }}
# - ternary: Ternary operator
#   {{ ternary "yes" "no" .condition }}
#
# Encoding:
# - b64enc, b64dec: Base64 encode/decode
# - b32enc, b32dec: Base32 encode/decode
#
# Lists:
# - list: Create list
#   {{ list 1 2 3 }}
# - first, rest, last, initial: List operations
# - append, prepend, concat: List manipulation
# - reverse, uniq, without, has: List utilities
# - slice: Slice list
#   {{ slice .list 1 3 }}
# - compact: Remove empty values
#
# Dictionaries:
# - dict: Create dictionary
#   {{ dict "key1" "val1" "key2" "val2" }}
# - get, set, unset: Dict operations (liner overrides these)
# - hasKey, pluck, keys, pick, omit, merge
# - deepCopy: Deep copy dict
# - mustDeepCopy: Deep copy with panic on error
#
# Type Conversion:
# - atoi: String to int
#   {{ "123" | atoi }}
# - int64, int, float64: Type conversions
# - toString: Convert to string
#
# Path & Filepath:
# - base, dir, ext, clean: Path operations
#   {{ "/path/to/file.txt" | base }}  // file.txt
# - isAbs: Check if absolute path
#
# Flow Control:
# - fail: Fail with error message
#   {{ if .error }}{{ fail "error occurred" }}{{ end }}
#
# UUID Functions:
# - uuidv4: Generate UUIDv4
#   {{ uuidv4 }}
#
# OS Functions:
# - env: Get environment variable
#   {{ env "HOME" }}
# - expandenv: Expand environment variables
#   {{ expandenv "$HOME/path" }}
#
# Semantic Version:
# - semver, semverCompare: Version comparison
#   {{ if semverCompare ">=1.2.3" .version }}ok{{ end }}
#
# Reflection:
# - typeOf, typeIs, typeIsLike, kindOf, kindIs
#   {{ typeOf .value }}
# - deepEqual: Deep equality check
#
# Cryptographic & Security:
# - sha1sum, sha256sum: Hash functions
#   {{ "text" | sha256sum }}
# - derivePassword: Derive password
# - genPrivateKey, buildCustomCert: Certificate generation
# - htpasswd: Generate htpasswd
#
# Network:
# - getHostByName: DNS lookup
#   {{ getHostByName "google.com" }}

## Liner Custom Functions:
#
# Dictionary Operations:
# - get(dict, key): Get value from dict
# - set(dict, key, value): Set value in dict
# - unset(dict, key): Remove key from dict
# - hasKey(dict, key): Check if key exists
#
# Network Functions:
# - host(hostport): Extract host from "host:port"
# - country(ip): Get country code
# - geoip(ip): Get GeoIP info (Country, City, ISP, ConnectionType)
# - geosite(domain): Get domain category
# - dnsResolve(host): Resolve to IP
# - domain(host): Extract root domain
# - hasIPv6(host): Check if host has IPv6
# - ipInt(ip): Convert IPv4 to uint32
# - ipRange(cidr): Parse CIDR range
# - isInNet(host, cidr): Check if in network
#
# String Matching:
# - hasPrefixes(pattern, str): Check prefixes (pipe-separated)
# - hasSuffixes(pattern, str): Check suffixes (pipe-separated)
# - regexMatch(pattern, str): Regex matching
# - wildcardMatch(pattern, str): Wildcard matching
#
# File Operations:
# - inFileLine(filename, line): Check line in file
# - inFileIPSet(filename, ip): Check IP in file
# - readFile(filename): Read file content
#
# HTTP:
# - fetch(ua, timeout, ttl, url): Fetch URL with cache
# - greased(clientHello): Check if TLS is greased
` + contextInfo)
}
