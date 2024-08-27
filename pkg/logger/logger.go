package logger

import (
    "io"
    "os"
    "runtime/debug"
    "strconv"
    "fmt"
    "time"
    "sync"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/pkgerrors"
    "gopkg.in/natefinch/lumberjack.v2"
)

const asciiArt = `
$$$$$$$\  $$$$$$\  $$$$$$\                         $$$$$\  $$$$$$\  $$\   $$\ $$\   $$\ 
$$  __$$\ \_$$  _|$$  __$$\                        \__$$ |$$  __$$\ $$ |  $$ |$$$\  $$ |
$$ |  $$ |  $$ |  $$ /  \__|                          $$ |$$ /  $$ |$$ |  $$ |$$$$\ $$ |
$$$$$$$\ |  $$ |  $$ |$$$$\       $$$$$$\             $$ |$$ |  $$ |$$$$$$$$ |$$ $$\$$ |
$$  __$$\   $$ |  $$ |\_$$ |      \______|      $$\   $$ |$$ |  $$ |$$  __$$ |$$ \$$$$ |
$$ |  $$ |  $$ |  $$ |  $$ |                    $$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |\$$$ |
$$$$$$$  |$$$$$$\ \$$$$$$  |                    \$$$$$$  | $$$$$$  |$$ |  $$ |$$ | \$$ |
\_______/ \______| \______/                      \______/  \______/ \__|  \__|\__|  \__|
                                                                                          
`

var (
    instance *Logger
    once     sync.Once
)

type Logger struct {
    zerolog.Logger
}

func PrintAsciiArt() {
    fmt.Print(asciiArt)
}

func Get() *Logger {
    once.Do(func() {
        instance = initLogger()
    })
    return instance
}

func initLogger() *Logger {
    zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
    zerolog.TimeFieldFormat = time.RFC3339Nano

    logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
    if err != nil {
        logLevel = int(zerolog.InfoLevel) // default to INFO
    }

    var output io.Writer = zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: time.RFC3339,
        FieldsExclude: []string{
            "user_agent",
            "git_revision",
            "go_version",
        },
    }

    if os.Getenv("APP_ENV") != "development" {
        fileLogger := &lumberjack.Logger{
            Filename:   "big-john-demo.log",
            MaxSize:    5,
            MaxBackups: 10,
            MaxAge:     14,
            Compress:   true,
        }

        output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
    }

    var gitRevision string

    buildInfo, ok := debug.ReadBuildInfo()
    if ok {
        for _, v := range buildInfo.Settings {
            if v.Key == "vcs.revision" {
                gitRevision = v.Value
                break
            }
        }
    }

    logger := zerolog.New(output).
        Level(zerolog.Level(logLevel)).
        With().
        Timestamp().
        Str("git_revision", gitRevision).
        Str("go_version", buildInfo.GoVersion).
        Caller().
        Logger()

    return &Logger{logger}
}

// Add convenience methods for different log levels
func (l *Logger) Info() *zerolog.Event {
    return l.Logger.Info()
}

func (l *Logger) Error() *zerolog.Event {
    return l.Logger.Error()
}

func (l *Logger) Warn() *zerolog.Event {
    return l.Logger.Warn()
}

func (l *Logger) Debug() *zerolog.Event {
    return l.Logger.Debug()
}

func (l *Logger) Fatal() *zerolog.Event {
    return l.Logger.Fatal()
}

