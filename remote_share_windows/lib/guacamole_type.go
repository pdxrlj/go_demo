package lib

const (
	defaultScreenWidth  = 1024
	defaultScreenHeight = 768
	defaultResolution   = 96
)

type GuacamoleConfig struct {
	ConnectionID string
	Protocol     string
	Parameters   map[string]string

	OptimalScreenWidth int

	OptimalScreenHeight int

	OptimalResolution int

	AudioMimetypes []string

	VideoMimetypes []string

	ImageMimetypes []string
}

type GuacamoleConfigOptions func(config *GuacamoleConfig)

func WithConnectionID(connectionId string) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.ConnectionID = connectionId
	}
}

func WithProtocol(protocol string) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.Protocol = protocol
	}
}

func WithParameters(parameters map[string]string) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.Parameters = parameters
	}
}

func WithOptimalScreenWidth(optimalScreenWidth int) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.OptimalScreenWidth = optimalScreenWidth
	}
}

func WithOptimalScreenHeight(optimalScreenHeight int) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.OptimalScreenHeight = optimalScreenHeight
	}
}

func WithOptimalResolution(optimalResolution int) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.OptimalResolution = optimalResolution
	}
}

func WithAudioMimetypes(audioMimetypes []string) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.AudioMimetypes = audioMimetypes
	}
}

func WithVideoMimetypes(videoMimetypes []string) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.VideoMimetypes = videoMimetypes
	}
}

func WithImageMimetypes(imageMimetypes []string) GuacamoleConfigOptions {
	return func(config *GuacamoleConfig) {
		config.ImageMimetypes = imageMimetypes
	}
}

func defaultGuacamoleConfig() *GuacamoleConfig {
	return &GuacamoleConfig{
		Parameters:          map[string]string{},
		OptimalScreenWidth:  defaultScreenWidth,
		OptimalScreenHeight: defaultScreenHeight,
		OptimalResolution:   defaultResolution,
		AudioMimetypes:      make([]string, 0, 1),
		VideoMimetypes:      make([]string, 0, 1),
		ImageMimetypes:      make([]string, 0, 1),
	}
}

func NewGuacamoleConfig(options ...GuacamoleConfigOptions) *GuacamoleConfig {
	d := defaultGuacamoleConfig()
	for _, option := range options {
		option(d)
	}
	return d
}
