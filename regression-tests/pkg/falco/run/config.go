package run

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RulesFile []string `yaml:"rules_file"`
	Plugins   []struct {
		Name        string      `yaml:"name"`
		LibraryPath string      `yaml:"library_path"`
		InitConfig  interface{} `yaml:"init_config,omitempty"` // todo: auto-marshal this to json
		OpenParams  string      `yaml:"open_params,omitempty"`
	} `yaml:"plugins"`
	LoadPlugins               []string `yaml:"load_plugins"`
	WatchConfigFiles          bool     `yaml:"watch_config_files"`
	TimeFormatIso8601         bool     `yaml:"time_format_iso_8601"`
	JSONOutput                bool     `yaml:"json_output"`
	JSONIncludeOutputProperty bool     `yaml:"json_include_output_property"`
	JSONIncludeTagsProperty   bool     `yaml:"json_include_tags_property"`
	LogStderr                 bool     `yaml:"log_stderr"`
	LogSyslog                 bool     `yaml:"log_syslog"`
	LogLevel                  string   `yaml:"log_level"`
	LibsLogger                struct {
		Enabled  bool   `yaml:"enabled"`
		Severity string `yaml:"severity"`
	} `yaml:"libs_logger"`
	Priority          string `yaml:"priority"`
	BufferedOutputs   bool   `yaml:"buffered_outputs"`
	SyscallEventDrops struct {
		Threshold     float64  `yaml:"threshold"`
		Actions       []string `yaml:"actions"`
		Rate          float64  `yaml:"rate"`
		MaxBurst      int      `yaml:"max_burst"`
		SimulateDrops bool     `yaml:"simulate_drops"`
	} `yaml:"syscall_event_drops"`
	SyscallEventTimeouts struct {
		MaxConsecutives int `yaml:"max_consecutives"`
	} `yaml:"syscall_event_timeouts"`
	SyscallBufSizePreset int `yaml:"syscall_buf_size_preset"`
	OutputTimeout        int `yaml:"output_timeout"`
	Outputs              struct {
		Rate     int `yaml:"rate"`
		MaxBurst int `yaml:"max_burst"`
	} `yaml:"outputs"`
	SyslogOutput struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"syslog_output"`
	FileOutput struct {
		Enabled   bool   `yaml:"enabled"`
		KeepAlive bool   `yaml:"keep_alive"`
		Filename  string `yaml:"filename"`
	} `yaml:"file_output"`
	StdoutOutput struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"stdout_output"`
	Webserver struct {
		Enabled            bool   `yaml:"enabled"`
		Threadiness        int    `yaml:"threadiness"`
		ListenPort         int    `yaml:"listen_port"`
		K8SHealthzEndpoint string `yaml:"k8s_healthz_endpoint"`
		SslEnabled         bool   `yaml:"ssl_enabled"`
		SslCertificate     string `yaml:"ssl_certificate"`
	} `yaml:"webserver"`
	ProgramOutput struct {
		Enabled   bool   `yaml:"enabled"`
		KeepAlive bool   `yaml:"keep_alive"`
		Program   string `yaml:"program"`
	} `yaml:"program_output"`
	HTTPOutput struct {
		Enabled   bool   `yaml:"enabled"`
		URL       string `yaml:"url"`
		UserAgent string `yaml:"user_agent"`
	} `yaml:"http_output"`
	Grpc struct {
		Enabled     bool   `yaml:"enabled"`
		BindAddress string `yaml:"bind_address"`
		Threadiness int    `yaml:"threadiness"`
	} `yaml:"grpc"`
	GrpcOutput struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"grpc_output"`
	MetadataDownload struct {
		MaxMb        int `yaml:"max_mb"`
		ChunkWaitUs  int `yaml:"chunk_wait_us"`
		WatchFreqSec int `yaml:"watch_freq_sec"`
	} `yaml:"metadata_download"`
}

func (c *Config) Marshal(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(c)
}

func (c *Config) Unmarshal(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}
