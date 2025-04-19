package config

import "time"

type Env struct {
	SystemENV
	FunctionalENV
	ProtocolENV
	//JWTENV
}

// SystemENV : 系統環境變數
type SystemENV struct {
	Environment            string        `env:"Environment" validate:"required"`
	Port                   int           `env:"Port" envDefault:"8080"`
	ReadTimeout            time.Duration `env:"ReadTimeout" envDefault:"1m0s"`
	WriteTimeout           time.Duration `env:"WriteTimeout" envDefault:"1m0s"`
	IdleTimeout            time.Duration `env:"IdleTimeout" envDefault:"5s"`
	FiberHeaderSizeLimitMb int           `env:"FiberHeaderSizeLimitMb" envDefault:"4"`
	FiberBodyLimitMb       int           `env:"FiberBodyLimitMb" envDefault:"4"`
	IsDebug                bool          `env:"IsDebug" envDefault:"false"`
	LogLevel               string        `env:"LogLevel" envDefault:"error"`
	LogType                string        `env:"LogType" envDefault:"console"`
	EnvMod                 string        `env:"EnvMod" envDefault:"dev"`
	RedisServer            string        `env:"REDIS_SERVER" envDefault:"localhost:6379"`
	RedisPassword          string        `env:"REDIS_PASSWORD" envDefault:"password"`
}

type FunctionalENV struct {
	Secure bool `env:"Secure" envDefault:"true"`
	MaxAge int  `env:"MaxAge" envDefault:"7200" validate:"min=1"`
	//GoogleClientID     string `env:"GoogleClientId" validate:"required"`
	//GoogleClientSecret string `env:"GoogleClientSecret" validate:"required"`
	//RefudAPIUrl        string `env:"RefudAPIUrl" validate:"required"`
}

type JWTENV struct {
	JwtSecretkey string `env:"jwtSecretkey" envDefault:"94EQD739HRW2443VC9986294G784D6351"`
	JwtExpires   int    `env:"jwtExpires" envDefault:"72000"`
	JwtTokenKey  string `env:"jwtTokenKey" envDefault:"jwtToken"`
	JwtUserKey   string `env:"jwtUserKey" envDefault:"currentUser"`
}

type ProtocolENV struct {
}

type LogLevel = string

const LogLevelInfo LogLevel = "info"
const LogLevelError LogLevel = "error"
