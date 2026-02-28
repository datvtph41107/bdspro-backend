package logger

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
)

type Config struct {
	Level       Level  `yaml:"level" default:"info"`
	Format      string `yaml:"format" default:"json"` // json, console
	OutputPath  string `yaml:"output_path" default:"stdout"`
	Development bool   `yaml:"development" default:"false"`
}
