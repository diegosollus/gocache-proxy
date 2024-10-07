package configs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// Define a estrutura para os recursos.
type Resource struct {
	Name           string `mapstructure:"name"`
	Endpoint       string `mapstructure:"endpoint"`
	DestinationURL string `mapstructure:"destination_url"`
}

// Define a estrutura principal da configuração.
type Configuration struct {
	Server struct {
		Host       string `mapstructure:"host"`
		ListenPort string `mapstructure:"listen_port"`
	} `mapstructure:"server"`
	Resources []Resource `mapstructure:"resources"`
}

var (
	config *Configuration
	once   sync.Once // Garante que o carregamento será feito uma única vez.
)

// Função que irá carregar a configuração.
func LoadConfig(path string) (*Configuration, error) {
	var err error

	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(path)

		// Suporte a variáveis de ambiente.
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Carrega o arquivo de configuração.
		if err = viper.ReadInConfig(); err != nil {
			err = fmt.Errorf("error loading config file: %w", err)
			return
		}

		// Unmarshal o conteúdo da configuração para a estrutura.
		if err = viper.Unmarshal(&config); err != nil {
			err = fmt.Errorf("error unmarshaling config file: %w", err)
		}
	})

	return config, err
}

// Função para obter a configuração carregada.
func GetConfig() *Configuration {
	return config
}
