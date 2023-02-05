package logger

type Format string

const (
	FormatJSON    Format = "json"
	FormatConsole Format = "console"
	sinkSTDOUT           = "stdout"
)

type FilterType string

type NameFilter struct {
	LoggerName string `yaml:"logger_name"`
}

// ExactSubnameFilter to filter out strings by name separated by dot.
// for "bar" given:
// * "b.bar.a" will filter out record.
// * "b.barFOO.a" will leave record be.
type ExactSubnameFilter struct {
	LoggerName string `yaml:"logger_name"`
}

type FilterConfig struct {
	FullNameFilter     []NameFilter         `yaml:"by_logger_name"`
	ExactSubnameFilter []ExactSubnameFilter `yaml:"by_exact_name"`
}

func DefaultFilterConfig() FilterConfig {
	return FilterConfig{}
}

type Config struct {
	Sink     string `yaml:"sink"`
	LogLevel string `yaml:"log_level"`
	Format   Format `yaml:"log_format"`
	// Filters allow filter out some log lines based on conditions
	Filters FilterConfig
}

func DefaultConfig() Config {
	return Config{
		Sink:     sinkSTDOUT,
		LogLevel: "DEBUG",
		Format:   FormatJSON,
		Filters:  DefaultFilterConfig(),
	}
}
