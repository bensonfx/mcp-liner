// Package config 提供liner配置结构定义和辅助函数
package config

import (
	"gopkg.in/yaml.v3"
)

// GlobalConfig 全局配置结构
type GlobalConfig struct {
	LogDir           string `json:"log_dir,omitempty" yaml:"log_dir,omitempty"`
	LogLevel         string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	LogBackups       int    `json:"log_backups,omitempty" yaml:"log_backups,omitempty"`
	LogMaxsize       int64  `json:"log_maxsize,omitempty" yaml:"log_maxsize,omitempty"`
	LogLocaltime     bool   `json:"log_localtime,omitempty" yaml:"log_localtime,omitempty"`
	LogChannelSize   uint   `json:"log_channel_size,omitempty" yaml:"log_channel_size,omitempty"`
	ForbidLocalAddr  bool   `json:"forbid_local_addr,omitempty" yaml:"forbid_local_addr,omitempty"`
	DialTimeout      int    `json:"dial_timeout,omitempty" yaml:"dial_timeout,omitempty"`
	DialReadBuffer   int    `json:"dial_read_buffer,omitempty" yaml:"dial_read_buffer,omitempty"`
	DialWriteBuffer  int    `json:"dial_write_buffer,omitempty" yaml:"dial_write_buffer,omitempty"`
	DnsServer        string `json:"dns_server,omitempty" yaml:"dns_server,omitempty"`
	DnsCacheDuration string `json:"dns_cache_duration,omitempty" yaml:"dns_cache_duration,omitempty"`
	DnsCacheSize     int    `json:"dns_cache_size,omitempty" yaml:"dns_cache_size,omitempty"`
	TcpReadBuffer    int    `json:"tcp_read_buffer,omitempty" yaml:"tcp_read_buffer,omitempty"`
	TcpWriteBuffer   int    `json:"tcp_write_buffer,omitempty" yaml:"tcp_write_buffer,omitempty"`
	TlsInsecure      bool   `json:"tls_insecure,omitempty" yaml:"tls_insecure,omitempty"`
	AutocertDir      string `json:"autocert_dir,omitempty" yaml:"autocert_dir,omitempty"`
	GeoipDir         string `json:"geoip_dir,omitempty" yaml:"geoip_dir,omitempty"`
	GeoipCacheSize   int    `json:"geoip_cache_size,omitempty" yaml:"geoip_cache_size,omitempty"`
	GeositeDisabled  bool   `json:"geosite_disabled,omitempty" yaml:"geosite_disabled,omitempty"`
	GeositeCacheSize int    `json:"geosite_cache_size,omitempty" yaml:"geosite_cache_size,omitempty"`
	IdleConnTimeout  int    `json:"idle_conn_timeout,omitempty" yaml:"idle_conn_timeout,omitempty"`
	MaxIdleConns     int    `json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`
	DisableHttp3     bool   `json:"disable_http3,omitempty" yaml:"disable_http3,omitempty"`
	SetProcessName   string `json:"set_process_name,omitempty" yaml:"set_process_name,omitempty"`
}

// HTTPForwardConfig HTTP转发配置
type HTTPForwardConfig struct {
	Policy           string `json:"policy,omitempty" yaml:"policy,omitempty"`
	AuthTable        string `json:"auth_table,omitempty" yaml:"auth_table,omitempty"`
	Dialer           string `json:"dialer,omitempty" yaml:"dialer,omitempty"`
	TcpCongestion    string `json:"tcp_congestion,omitempty" yaml:"tcp_congestion,omitempty"`
	DenyDomainsTable string `json:"deny_domains_table,omitempty" yaml:"deny_domains_table,omitempty"`
	SpeedLimit       int64  `json:"speed_limit,omitempty" yaml:"speed_limit,omitempty"`
	DisableIpv6      bool   `json:"disable_ipv6,omitempty" yaml:"disable_ipv6,omitempty"`
	PreferIpv6       bool   `json:"prefer_ipv6,omitempty" yaml:"prefer_ipv6,omitempty"`
	Log              bool   `json:"log,omitempty" yaml:"log,omitempty"`
	LogInterval      int64  `json:"log_interval,omitempty" yaml:"log_interval,omitempty"`
	IoCopyBuffer     int    `json:"io_copy_buffer,omitempty" yaml:"io_copy_buffer,omitempty"`
	IdleTimeout      int64  `json:"idle_timeout,omitempty" yaml:"idle_timeout,omitempty"`
}

// HTTPTunnelConfig HTTP隧道配置
type HTTPTunnelConfig struct {
	Enabled         bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	AuthTable       string   `json:"auth_table,omitempty" yaml:"auth_table,omitempty"`
	AllowListens    []string `json:"allow_listens,omitempty" yaml:"allow_listens,omitempty"`
	SpeedLimit      int64    `json:"speed_limit,omitempty" yaml:"speed_limit,omitempty"`
	EnableKeepAlive bool     `json:"enable_keep_alive,omitempty" yaml:"enable_keep_alive,omitempty"`
	Log             bool     `json:"log,omitempty" yaml:"log,omitempty"`
}

// HTTPWebIndexConfig Web Index配置
type HTTPWebIndexConfig struct {
	Root    string `json:"root,omitempty" yaml:"root,omitempty"`
	Headers string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Charset string `json:"charset,omitempty" yaml:"charset,omitempty"`
	Body    string `json:"body,omitempty" yaml:"body,omitempty"`
	File    string `json:"file,omitempty" yaml:"file,omitempty"`
}

// HTTPWebProxyConfig Web Proxy配置
type HTTPWebProxyConfig struct {
	Pass        string `json:"pass,omitempty" yaml:"pass,omitempty"`
	AuthTable   string `json:"auth_table,omitempty" yaml:"auth_table,omitempty"`
	StripPrefix string `json:"strip_prefix,omitempty" yaml:"strip_prefix,omitempty"`
	SetHeaders  string `json:"set_headers,omitempty" yaml:"set_headers,omitempty"`
	DumpFailure bool   `json:"dump_failure,omitempty" yaml:"dump_failure,omitempty"`
}

// HTTPWebDohConfig Web DoH配置
type HTTPWebDohConfig struct {
	Enabled   bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Policy    string `json:"policy,omitempty" yaml:"policy,omitempty"`
	ProxyPass string `json:"proxy_pass,omitempty" yaml:"proxy_pass,omitempty"`
	CacheSize int    `json:"cache_size,omitempty" yaml:"cache_size,omitempty"`
}

// HTTPWebConfig Web配置
type HTTPWebConfig struct {
	Location      string             `json:"location,omitempty" yaml:"location,omitempty"`
	TcpCongestion string             `json:"tcp_congestion,omitempty" yaml:"tcp_congestion,omitempty"`
	Index         HTTPWebIndexConfig `json:"index,omitempty" yaml:"index,omitempty"`
	Proxy         HTTPWebProxyConfig `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	Doh           HTTPWebDohConfig   `json:"doh,omitempty" yaml:"doh,omitempty"`
}

// HTTPConfig HTTP/HTTPS配置
type HTTPConfig struct {
	Listen     []string          `json:"listen,omitempty" yaml:"listen,omitempty"`
	ServerName []string          `json:"server_name,omitempty" yaml:"server_name,omitempty"`
	Keyfile    string            `json:"keyfile,omitempty" yaml:"keyfile,omitempty"`
	Certfile   string            `json:"certfile,omitempty" yaml:"certfile,omitempty"`
	PSK        string            `json:"psk,omitempty" yaml:"psk,omitempty"`
	Forward    HTTPForwardConfig `json:"forward,omitempty" yaml:"forward,omitempty"`
	Tunnel     HTTPTunnelConfig  `json:"tunnel,omitempty" yaml:"tunnel,omitempty"`
	Web        []HTTPWebConfig   `json:"web,omitempty" yaml:"web,omitempty"`
}

// TunnelConfig 隧道配置
type TunnelConfig struct {
	RemoteListen    []string `json:"remote_listen,omitempty" yaml:"remote_listen,omitempty"`
	ProxyPass       string   `json:"proxy_pass,omitempty" yaml:"proxy_pass,omitempty"`
	Resolver        string   `json:"resolver,omitempty" yaml:"resolver,omitempty"`
	DialTimeout     int      `json:"dial_timeout,omitempty" yaml:"dial_timeout,omitempty"`
	Dialer          string   `json:"dialer,omitempty" yaml:"dialer,omitempty"`
	SpeedLimit      int64    `json:"speed_limit,omitempty" yaml:"speed_limit,omitempty"`
	EnableKeepAlive bool     `json:"enable_keep_alive,omitempty" yaml:"enable_keep_alive,omitempty"`
	Log             bool     `json:"log,omitempty" yaml:"log,omitempty"`
}

// DnsConfig DNS配置
type DnsConfig struct {
	Listen    []string `json:"listen,omitempty" yaml:"listen,omitempty"`
	Keyfile   string   `json:"keyfile,omitempty" yaml:"keyfile,omitempty"`
	Policy    string   `json:"policy,omitempty" yaml:"policy,omitempty"`
	ProxyPass string   `json:"proxy_pass,omitempty" yaml:"proxy_pass,omitempty"`
	CacheSize int      `json:"cache_size,omitempty" yaml:"cache_size,omitempty"`
	Log       bool     `json:"log,omitempty" yaml:"log,omitempty"`
}

// SocksForwardConfig Socks转发配置
type SocksForwardConfig struct {
	Policy           string `json:"policy,omitempty" yaml:"policy,omitempty"`
	AuthTable        string `json:"auth_table,omitempty" yaml:"auth_table,omitempty"`
	Dialer           string `json:"dialer,omitempty" yaml:"dialer,omitempty"`
	DenyDomainsTable string `json:"deny_domains_table,omitempty" yaml:"deny_domains_table,omitempty"`
	SpeedLimit       int64  `json:"speed_limit,omitempty" yaml:"speed_limit,omitempty"`
	DisableIpv6      bool   `json:"disable_ipv6,omitempty" yaml:"disable_ipv6,omitempty"`
	PreferIpv6       bool   `json:"prefer_ipv6,omitempty" yaml:"prefer_ipv6,omitempty"`
	Log              bool   `json:"log,omitempty" yaml:"log,omitempty"`
}

// SocksConfig Socks代理配置
type SocksConfig struct {
	Listen  []string           `json:"listen,omitempty" yaml:"listen,omitempty"`
	PSK     string             `json:"psk,omitempty" yaml:"psk,omitempty"`
	Forward SocksForwardConfig `json:"forward,omitempty" yaml:"forward,omitempty"`
}

// SniForwardConfig SNI转发配置
type SniForwardConfig struct {
	Policy      string `json:"policy,omitempty" yaml:"policy,omitempty"`
	Dialer      string `json:"dialer,omitempty" yaml:"dialer,omitempty"`
	DisableIpv6 bool   `json:"disable_ipv6,omitempty" yaml:"disable_ipv6,omitempty"`
	PreferIpv6  bool   `json:"prefer_ipv6,omitempty" yaml:"prefer_ipv6,omitempty"`
	Log         bool   `json:"log,omitempty" yaml:"log,omitempty"`
}

// SniConfig SNI配置
type SniConfig struct {
	Enabled bool             `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Forward SniForwardConfig `json:"forward,omitempty" yaml:"forward,omitempty"`
}

// RedsocksForwardConfig Redsocks转发配置
type RedsocksForwardConfig struct {
	Dialer string `json:"dialer,omitempty" yaml:"dialer,omitempty"`
	Log    bool   `json:"log,omitempty" yaml:"log,omitempty"`
}

// RedsocksConfig Redsocks透明代理配置
type RedsocksConfig struct {
	Listen  []string              `json:"listen,omitempty" yaml:"listen,omitempty"`
	Forward RedsocksForwardConfig `json:"forward,omitempty" yaml:"forward,omitempty"`
}

// StreamConfig 流转发配置
type StreamConfig struct {
	Listen        []string `json:"listen,omitempty" yaml:"listen,omitempty"`
	Keyfile       string   `json:"keyfile,omitempty" yaml:"keyfile,omitempty"`
	Certfile      string   `json:"certfile,omitempty" yaml:"certfile,omitempty"`
	ProxyPass     string   `json:"proxy_pass,omitempty" yaml:"proxy_pass,omitempty"`
	ProxyProtocol uint     `json:"proxy_protocol,omitempty" yaml:"proxy_protocol,omitempty"`
	DialTimeout   int      `json:"dial_timeout,omitempty" yaml:"dial_timeout,omitempty"`
	Dialer        string   `json:"dialer,omitempty" yaml:"dialer,omitempty"`
	SpeedLimit    int64    `json:"speed_limit,omitempty" yaml:"speed_limit,omitempty"`
	Log           bool     `json:"log,omitempty" yaml:"log,omitempty"`
}

// Config liner完整配置
type Config struct {
	Global   GlobalConfig      `json:"global,omitempty" yaml:"global,omitempty"`
	Dialer   map[string]string `json:"dialer,omitempty" yaml:"dialer,omitempty"`
	Sni      SniConfig         `json:"sni,omitempty" yaml:"sni,omitempty"`
	Https    []HTTPConfig      `json:"https,omitempty" yaml:"https,omitempty"`
	Http     []HTTPConfig      `json:"http,omitempty" yaml:"http,omitempty"`
	Tunnel   []TunnelConfig    `json:"tunnel,omitempty" yaml:"tunnel,omitempty"`
	Dns      []DnsConfig       `json:"dns,omitempty" yaml:"dns,omitempty"`
	Socks    []SocksConfig     `json:"socks,omitempty" yaml:"socks,omitempty"`
	Redsocks []RedsocksConfig  `json:"redsocks,omitempty" yaml:"redsocks,omitempty"`
	Stream   []StreamConfig    `json:"stream,omitempty" yaml:"stream,omitempty"`
}

// NewDefaultGlobalConfig 创建默认全局配置
func NewDefaultGlobalConfig() GlobalConfig {
	return GlobalConfig{
		LogLevel:        "info",
		DnsServer:       "https://8.8.8.8/dns-query",
		DisableHttp3:    false,
		DialTimeout:     5,
		IdleConnTimeout: 90,
		MaxIdleConns:    100,
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
