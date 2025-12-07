// Package config 提供liner配置结构定义和辅助函数
package config

import (
	"gopkg.in/yaml.v3"
)

// GlobalConfig 全局配置结构，完全对应 liner/config.go Config.Global
type GlobalConfig struct {
	LogDir           string `json:"log_dir" yaml:"log_dir"`
	LogLevel         string `json:"log_level" yaml:"log_level"`
	LogBackups       int    `json:"log_backups" yaml:"log_backups"`
	LogMaxsize       int64  `json:"log_maxsize" yaml:"log_maxsize"`
	LogLocaltime     bool   `json:"log_localtime" yaml:"log_localtime"`
	LogChannelSize   uint   `json:"log_channel_size" yaml:"log_channel_size"`
	ForbidLocalAddr  bool   `json:"forbid_local_addr" yaml:"forbid_local_addr"`
	DialTimeout      int    `json:"dial_timeout" yaml:"dial_timeout"`
	DialReadBuffer   int    `json:"dial_read_buffer" yaml:"dial_read_buffer"` // Danger, see https://issues.apache.org/jira/browse/KAFKA-16496
	DialWriteBuffer  int    `json:"dial_write_buffer" yaml:"dial_write_buffer"`
	DnsServer        string `json:"dns_server" yaml:"dns_server"`
	DnsCacheDuration string `json:"dns_cache_duration" yaml:"dns_cache_duration"`
	DnsCacheSize     int    `json:"dns_cache_size" yaml:"dns_cache_size"`
	TcpReadBuffer    int    `json:"tcp_read_buffer" yaml:"tcp_read_buffer"`
	TcpWriteBuffer   int    `json:"tcp_write_buffer" yaml:"tcp_write_buffer"`
	TlsInsecure      bool   `json:"tls_insecure" yaml:"tls_insecure"`
	AutocertDir      string `json:"autocert_dir" yaml:"autocert_dir"`
	GeoipDir         string `json:"geoip_dir" yaml:"geoip_dir"`
	GeoipCacheSize   int    `json:"geoip_cache_size" yaml:"geoip_cache_size"`
	GeositeDisabled  bool   `json:"geosite_disabled" yaml:"geosite_disabled"`
	GeositeCacheSize int    `json:"geosite_cache_size" yaml:"geosite_cache_size"`
	IdleConnTimeout  int    `json:"idle_conn_timeout" yaml:"idle_conn_timeout"`
	MaxIdleConns     int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	DisableHttp3     bool   `json:"disable_http3" yaml:"disable_http3"`
	SetProcessName   string `json:"set_process_name" yaml:"set_process_name"`
}

// HTTPForwardConfig HTTP转发配置，对应 liner/config.go HTTPConfig.Forward
type HTTPForwardConfig struct {
	Policy           string `json:"policy" yaml:"policy"`
	AuthTable        string `json:"auth_table" yaml:"auth_table"`
	Dialer           string `json:"dialer" yaml:"dialer"`
	TcpCongestion    string `json:"tcp_congestion" yaml:"tcp_congestion"`
	DenyDomainsTable string `json:"deny_domains_table" yaml:"deny_domains_table"`
	SpeedLimit       int64  `json:"speed_limit" yaml:"speed_limit"`
	DisableIpv6      bool   `json:"disable_ipv6" yaml:"disable_ipv6"`
	PreferIpv6       bool   `json:"prefer_ipv6" yaml:"prefer_ipv6"`
	Log              bool   `json:"log" yaml:"log"`
	LogInterval      int64  `json:"log_interval" yaml:"log_interval"`
	IoCopyBuffer     int    `json:"io_copy_buffer" yaml:"io_copy_buffer"`
	IdleTimeout      int64  `json:"idle_timeout" yaml:"idle_timeout"`
}

// HTTPTunnelConfig HTTP隧道配置，对应 liner/config.go HTTPConfig.Tunnel
type HTTPTunnelConfig struct {
	Enabled         bool     `json:"enabled" yaml:"enabled"`
	AuthTable       string   `json:"auth_table" yaml:"auth_table"`
	AllowListens    []string `json:"allow_listens" yaml:"allow_listens"`
	SpeedLimit      int64    `json:"speed_limit" yaml:"speed_limit"`
	EnableKeepAlive bool     `json:"enable_keep_alive" yaml:"enable_keep_alive"`
	Log             bool     `json:"log" yaml:"log"`
}

// HTTPWebIndexConfig Web Index配置
type HTTPWebIndexConfig struct {
	Root    string `json:"root" yaml:"root"`
	Headers string `json:"headers" yaml:"headers"`
	Charset string `json:"charset" yaml:"charset"`
	Body    string `json:"body" yaml:"body"`
	File    string `json:"file" yaml:"file"`
}

// HTTPWebProxyConfig Web Proxy配置
type HTTPWebProxyConfig struct {
	Pass        string `json:"pass" yaml:"pass"`
	AuthTable   string `json:"auth_table" yaml:"auth_table"`
	StripPrefix string `json:"strip_prefix" yaml:"strip_prefix"`
	SetHeaders  string `json:"set_headers" yaml:"set_headers"`
	DumpFailure bool   `json:"dump_failure" yaml:"dump_failure"`
}

// HTTPWebDohConfig Web DoH配置
type HTTPWebDohConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Policy    string `json:"policy" yaml:"policy"`
	ProxyPass string `json:"proxy_pass" yaml:"proxy_pass"`
	CacheSize int    `json:"cache_size" yaml:"cache_size"`
}

// HTTPWebFastcgiConfig FastCGI配置
type HTTPWebFastcgiConfig struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Root       string `json:"root" yaml:"root"`
	DefaultAPP string `json:"default_app" yaml:"default_app"`
}

// HTTPWebDavConfig WebDAV配置
type HTTPWebDavConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Root      string `json:"root" yaml:"root"`
	AuthTable string `json:"auth_table" yaml:"auth_table"`
}

// HTTPWebShellConfig Web Shell配置
type HTTPWebShellConfig struct {
	Enabled   bool              `json:"enabled" yaml:"enabled"`
	AuthTable string            `json:"auth_table" yaml:"auth_table"`
	Command   string            `json:"command" yaml:"command"`
	Home      string            `json:"home" yaml:"home"`
	Template  map[string]string `json:"template" yaml:"template"`
}

// HTTPWebConfig 对应 liner/config.go HTTPConfig.Web
type HTTPWebConfig struct {
	Location      string               `json:"location" yaml:"location"`
	TcpCongestion string               `json:"tcp_congestion" yaml:"tcp_congestion"`
	Fastcgi       HTTPWebFastcgiConfig `json:"fastcgi" yaml:"fastcgi"`
	Dav           HTTPWebDavConfig     `json:"dav" yaml:"dav"`
	Doh           HTTPWebDohConfig     `json:"doh" yaml:"doh"`
	Index         HTTPWebIndexConfig   `json:"index" yaml:"index"`
	Proxy         HTTPWebProxyConfig   `json:"proxy" yaml:"proxy"`
	Shell         HTTPWebShellConfig   `json:"shell" yaml:"shell"`
}

// ServerConfig 对应 liner/config.go HTTPConfig.ServerConfig 中的值
type ServerConfig struct {
	Keyfile        string `json:"keyfile" yaml:"keyfile"`
	Certfile       string `json:"certfile" yaml:"certfile"`
	DisableHttp2   bool   `json:"disable_http2" yaml:"disable_http2"`
	DisableHttp3   bool   `json:"disable_http3" yaml:"disable_http3"`
	DisableTls11   bool   `json:"disable_tls11" yaml:"disable_tls11"`
	DisableOcsp    bool   `json:"disable_ocsp" yaml:"disable_ocsp"`
	PreferChacha20 bool   `json:"prefer_chacha20" yaml:"prefer_chacha20"`
}

// HTTPConfig HTTP/HTTPS配置，对应 liner/config.go HTTPConfig
type HTTPConfig struct {
	Listen       []string                `json:"listen" yaml:"listen"`
	ServerName   []string                `json:"server_name" yaml:"server_name"`
	Keyfile      string                  `json:"keyfile" yaml:"keyfile"`
	Certfile     string                  `json:"certfile" yaml:"certfile"`
	PSK          string                  `json:"psk" yaml:"psk"`
	ServerConfig map[string]ServerConfig `json:"server_config" yaml:"server_config"`
	Forward      HTTPForwardConfig       `json:"forward" yaml:"forward"`
	Tunnel       HTTPTunnelConfig        `json:"tunnel" yaml:"tunnel"`
	Web          []HTTPWebConfig         `json:"web" yaml:"web"`
}

// TunnelConfig 隧道配置，对应 liner/config.go TunnelConfig
type TunnelConfig struct {
	RemoteListen    []string `json:"remote_listen" yaml:"remote_listen"`
	ProxyPass       string   `json:"proxy_pass" yaml:"proxy_pass"`
	Resolver        string   `json:"resolver" yaml:"resolver"`
	DialTimeout     int      `json:"dial_timeout" yaml:"dial_timeout"`
	Dialer          string   `json:"dialer" yaml:"dialer"`
	SpeedLimit      int64    `json:"speed_limit" yaml:"speed_limit"`
	EnableKeepAlive bool     `json:"enable_keep_alive" yaml:"enable_keep_alive"`
	Log             bool     `json:"log" yaml:"log"`
}

// DnsConfig DNS配置，对应 liner/config.go DnsConfig
type DnsConfig struct {
	Listen    []string `json:"listen" yaml:"listen"`
	Keyfile   string   `json:"keyfile" yaml:"keyfile"`
	Policy    string   `json:"policy" yaml:"policy"`
	ProxyPass string   `json:"proxy_pass" yaml:"proxy_pass"`
	CacheSize int      `json:"cache_size" yaml:"cache_size"`
	Log       bool     `json:"log" yaml:"log"`
}

// SocksForwardConfig Socks转发配置
type SocksForwardConfig struct {
	Policy           string `json:"policy" yaml:"policy"`
	AuthTable        string `json:"auth_table" yaml:"auth_table"`
	Dialer           string `json:"dialer" yaml:"dialer"`
	DenyDomainsTable string `json:"deny_domains_table" yaml:"deny_domains_table"`
	SpeedLimit       int64  `json:"speed_limit" yaml:"speed_limit"`
	DisableIpv6      bool   `json:"disable_ipv6" yaml:"disable_ipv6"`
	PreferIpv6       bool   `json:"prefer_ipv6" yaml:"prefer_ipv6"`
	Log              bool   `json:"log" yaml:"log"`
}

// SocksConfig Socks代理配置，对应 liner/config.go SocksConfig
type SocksConfig struct {
	Listen  []string           `json:"listen" yaml:"listen"`
	PSK     string             `json:"psk" yaml:"psk"`
	Forward SocksForwardConfig `json:"forward" yaml:"forward"`
}

// SniForwardConfig SNI转发配置
type SniForwardConfig struct {
	Policy      string `json:"policy" yaml:"policy"`
	Dialer      string `json:"dialer" yaml:"dialer"`
	DisableIpv6 bool   `json:"disable_ipv6" yaml:"disable_ipv6"`
	PreferIpv6  bool   `json:"prefer_ipv6" yaml:"prefer_ipv6"`
	Log         bool   `json:"log" yaml:"log"`
}

// SniConfig SNI配置 (对应 liner/config.go SniConfig)
type SniConfig struct {
	Enabled bool             `json:"enabled" yaml:"enabled"`
	Forward SniForwardConfig `json:"forward" yaml:"forward"`
}

// RedsocksForwardConfig Redsocks转发配置
type RedsocksForwardConfig struct {
	Dialer string `json:"dialer" yaml:"dialer"`
	Log    bool   `json:"log" yaml:"log"`
}

// RedsocksConfig Redsocks配置 (对应 liner/config.go RedsocksConfig)
type RedsocksConfig struct {
	Listen  []string              `json:"listen" yaml:"listen"`
	Forward RedsocksForwardConfig `json:"forward" yaml:"forward"`
}

// StreamConfig 流转发配置 (对应 liner/config.go StreamConfig)
type StreamConfig struct {
	Listen        []string `json:"listen" yaml:"listen"`
	Keyfile       string   `json:"keyfile" yaml:"keyfile"`
	Certfile      string   `json:"certfile" yaml:"certfile"`
	ProxyPass     string   `json:"proxy_pass" yaml:"proxy_pass"`
	ProxyProtocol uint     `json:"proxy_protocol" yaml:"proxy_protocol"`
	DialTimeout   int      `json:"dial_timeout" yaml:"dial_timeout"`
	Dialer        string   `json:"dialer" yaml:"dialer"`
	SpeedLimit    int64    `json:"speed_limit" yaml:"speed_limit"`
	Log           bool     `json:"log" yaml:"log"`
}

// SshConfig SSH配置 (对应 liner/config.go SshConfig)
type SshConfig struct {
	Listen           []string `json:"listen" yaml:"listen"`
	ServerVersion    string   `json:"server_version" yaml:"server_version"`
	TcpReadBuffer    int      `json:"tcp_read_buffer" yaml:"tcp_read_buffer"`
	TcpWriteBuffer   int      `json:"tcp_write_buffer" yaml:"tcp_write_buffer"`
	DisableKeepalive bool     `json:"disable_keepalive" yaml:"disable_keepalive"`
	BannerFile       string   `json:"banner_file" yaml:"banner_file"`
	HostKey          string   `json:"host_key" yaml:"host_key"`
	AuthTable        string   `json:"auth_table" yaml:"auth_table"`
	AuthorizedKeys   string   `json:"authorized_keys" yaml:"authorized_keys"`
	Shell            string   `json:"shell" yaml:"shell"`
	Home             string   `json:"home" yaml:"home"`
	EnvFile          string   `json:"env_file" yaml:"env_file"`
	Log              bool     `json:"log" yaml:"log"`
}

// CronConfig (对应 liner/config.go Config.Cron)
type CronConfig struct {
	Spec    string `json:"spec" yaml:"spec"`
	Command string `json:"command" yaml:"command"`
}

// Config liner完整配置结构，对应 liner/config.go Config
type Config struct {
	Global   GlobalConfig      `json:"global" yaml:"global"`
	Cron     []CronConfig      `json:"cron" yaml:"cron"`
	Dialer   map[string]string `json:"dialer" yaml:"dialer"`
	Sni      SniConfig         `json:"sni" yaml:"sni"`
	Https    []HTTPConfig      `json:"https" yaml:"https"`
	Http     []HTTPConfig      `json:"http" yaml:"http"`
	Socks    []SocksConfig     `json:"socks" yaml:"socks"`
	Redsocks []RedsocksConfig  `json:"redsocks" yaml:"redsocks"`
	Tunnel   []TunnelConfig    `json:"tunnel" yaml:"tunnel"`
	Stream   []StreamConfig    `json:"stream" yaml:"stream"`
	Ssh      []SshConfig       `json:"ssh" yaml:"ssh"`
	Dns      []DnsConfig       `json:"dns" yaml:"dns"`
}

// NewDefaultGlobalConfig 创建默认全局配置
// 默认值与 liner/example.yaml 保持一致
func NewDefaultGlobalConfig() GlobalConfig {
	return GlobalConfig{
		LogLevel:         "info",
		LogBackups:       2,
		LogMaxsize:       1073741824, // 1GB
		LogLocaltime:     true,
		MaxIdleConns:     100,
		DialTimeout:      30,
		DnsCacheDuration: "15m",
		DnsCacheSize:     524288,
		DnsServer:        "https://8.8.8.8/dns-query",
	}
}

// NewDefaultHTTPConfig 创建默认HTTP配置
func NewDefaultHTTPConfig(listen []string, serverName []string) HTTPConfig {
	return HTTPConfig{
		Listen:     listen,
		ServerName: serverName,
		Forward: HTTPForwardConfig{
			Policy: "proxy_pass",
			Dialer: "local",
			Log:    true,
		},
	}
}

// NewDefaultTunnelConfig 创建默认隧道配置
func NewDefaultTunnelConfig(remoteListen []string, proxyPass string, dialer string) TunnelConfig {
	return TunnelConfig{
		RemoteListen:    remoteListen,
		ProxyPass:       proxyPass,
		Dialer:          dialer,
		DialTimeout:     5,
		EnableKeepAlive: true,
		Log:             true,
	}
}

// NewDefaultDNSConfig 创建默认DNS配置
func NewDefaultDNSConfig(listen []string, proxyPass string) DnsConfig {
	return DnsConfig{
		Listen:    listen,
		Policy:    "forward",
		ProxyPass: proxyPass,
		CacheSize: 4096,
		Log:       true,
	}
}

// ToYAML 将配置序列化为YAML字符串
func (c *Config) ToYAML() (string, error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromYAML 从YAML字符串解析配置
func FromYAML(yamlStr string) (*Config, error) {
	var config Config
	err := yaml.Unmarshal([]byte(yamlStr), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
