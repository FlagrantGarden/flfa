package flfa

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	ApiVersionKey      string = "apiVersion"
	ApiVersionDefault  string = "0.1.0"
	ModulePathKey      string = "modulePath"
	UserDataPathKey    string = "userDataPath"
	CurrentGameKey     string = "currentGame"
	CurrentGameDefault string = "default"
)

type Config struct {
	ApiVersion   *semver.Version
	ModulePath   string
	UserDataPath string
	CurrentGame  string
}

func (f *Api) GenerateDefaultConfig() {
	parsedApiVersion := f.GetDefaultApiVersion()
	log.Trace().Msgf("setting default config (%s:%s)", ApiVersionKey, parsedApiVersion)
	viper.GetViper().SetDefault(ApiVersionKey, parsedApiVersion)

	defaultModulePath, err := f.GetDefaultModulePath()
	if err != nil {
		log.Panic().Msgf("unable to generate default config value for '%s': %s", ModulePathKey, err)
	}
	log.Trace().Msgf("setting default module path (%s: %s)", ModulePathKey, defaultModulePath)
	viper.SetDefault(ModulePathKey, defaultModulePath)

	defaultUserDataPath, err := f.GetDefaultUserDataPath()
	if err != nil {
		log.Panic().Msgf("unable to generate default config value for '%s': %s", UserDataPathKey, err)
	}
	log.Trace().Msgf("setting default user data path (%s: %s)", UserDataPathKey, defaultUserDataPath)
	viper.SetDefault(UserDataPathKey, defaultUserDataPath)
}

func (f *Api) LoadConfig() (err error) {
	var semanticApiVersion *semver.Version
	configuredApiVersion := viper.GetString(ApiVersionKey)
	if configuredApiVersion == "" {
		semanticApiVersion = f.GetDefaultApiVersion()
	} else {
		semanticApiVersion, err = semver.NewVersion(configuredApiVersion)
		if err != nil {
			return fmt.Errorf("could not load '%s' from config '%s': %s", ApiVersionKey, viper.GetViper().ConfigFileUsed(), err)
		}
	}
	f.RunningConfig.ApiVersion = semanticApiVersion

	f.RunningConfig.ModulePath = viper.GetString(ModulePathKey)

	f.RunningConfig.UserDataPath = viper.GetString(UserDataPathKey)

	f.RunningConfig.CurrentGame = viper.GetString(CurrentGameKey)

	return nil
}

func (f *Api) GetDefaultModulePath() (string, error) {
	executableDirectory, err := os.Executable()
	if err != nil {
		return "", err
	}

	defaultModulePath := filepath.Join(filepath.Dir(executableDirectory), "modules")
	log.Trace().Msgf("default config for module path: %s", defaultModulePath)
	return defaultModulePath, nil
}

func (f *Api) GetDefaultUserDataPath() (string, error) {
	executableDirectory, err := os.Executable()
	if err != nil {
		return "", err
	}

	defaultUserDataPath := filepath.Join(filepath.Dir(executableDirectory), "data")
	log.Trace().Msgf("default config for module path: %s", defaultUserDataPath)
	return defaultUserDataPath, nil
}

func (f *Api) GetDefaultApiVersion() *semver.Version {
	parsedApiVersion, err := semver.NewVersion(ApiVersionDefault)
	if err != nil {
		log.Panic().Msgf("unable to generate default config value for '%s': %s", ApiVersionKey, err)
	}
	return parsedApiVersion
}
