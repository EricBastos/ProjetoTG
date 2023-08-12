package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

var Cfg *Conf

type JwtConf struct {
	JwtSecret     string `mapstructure:"JWT_SECRET"`
	JwtExpiration int    `mapstructure:"JWT_EXPIRATION"`
	TokenAuthUser *jwtauth.JWTAuth
}

type DBConf struct {
	DBHost     string `mapstructure:"HOST_DB"`
	DBPort     string `mapstructure:"PORT_DB"`
	DBName     string `mapstructure:"DATABASE_NAME"`
	DBUser     string `mapstructure:"USER_DB"`
	DBPassword string `mapstructure:"PASS_DB"`
}

type BankWebhookConf struct {
	BankWebhookHost string `mapstructure:"BANK_WEBHOOK_HOST"`
	BankWebhookPort string `mapstructure:"BANK_WEBHOOK_PORT"`
}

type RabbitmqConf struct {
	RABBITHost     string `mapstructure:"RABBITMQ_HOST"`
	RABBITPort     string `mapstructure:"RABBITMQ_PORT"`
	RABBITUser     string `mapstructure:"RABBITMQ_USER"`
	RABBITPassword string `mapstructure:"RABBITMQ_PASS"`

	RABBITCallExchange     string `mapstructure:"RABBITMQ_CALL_EXCHANGE"`
	RABBITFeedbackExchange string `mapstructure:"RABBITMQ_FEEDBACK_EXCHANGE"`

	RABBITCallQueueEthereum string `mapstructure:"RABBITMQ_CALL_QUEUE_ETHEREUM"`
	RABBITCallQueuePolygon  string `mapstructure:"RABBITMQ_CALL_QUEUE_POLYGON"`
}

type WalletsConf struct {
	EthereumWalletAddress    string
	EthereumWalletPrivateKey string
	PolygonWalletAddress     string
	PolygonWalletPrivateKey  string
}

type ContractsConf struct {
	EthereumRpcHost       string
	EthereumTokenContract string
	PolygonRpcHost        string
	PolygonTokenContract  string
}

type Conf struct {
	DBConf
	WalletsConf
	ContractsConf
	RabbitmqConf
	JwtConf
	BankWebhookConf
}

func LoadConfig(path string) (*Conf, error) {
	var dbCfg DBConf
	var rabbitCfg RabbitmqConf
	var jwtCfg JwtConf
	var bankWebhookConf BankWebhookConf
	var walletsConf WalletsConf
	var contractsConf ContractsConf
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&dbCfg)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&rabbitCfg)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&jwtCfg)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&bankWebhookConf)
	if err != nil {
		return nil, err
	}

	walletsConf.EthereumWalletAddress = viper.GetString("ETH_PUBLIC_KEY")
	walletsConf.EthereumWalletPrivateKey = viper.GetString("ETH_PRIVATE_KEY")
	walletsConf.PolygonWalletAddress = viper.GetString("POLYGON_PUBLIC_KEY")
	walletsConf.PolygonWalletPrivateKey = viper.GetString("POLYGON_PRIVATE_KEY")

	contractsConf.EthereumRpcHost = viper.GetString("ETH_RPC_HOST")
	contractsConf.EthereumTokenContract = viper.GetString("ETH_CONTRACT_ADDRESS")
	contractsConf.PolygonRpcHost = viper.GetString("POLYGON_RPC_HOST")
	contractsConf.PolygonTokenContract = viper.GetString("POLYGON_CONTRACT_ADDRESS")

	cfg := &Conf{
		DBConf:          dbCfg,
		RabbitmqConf:    rabbitCfg,
		JwtConf:         jwtCfg,
		WalletsConf:     walletsConf,
		ContractsConf:   contractsConf,
		BankWebhookConf: bankWebhookConf,
	}
	cfg.TokenAuthUser = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)
	Cfg = cfg
	return cfg, err
}