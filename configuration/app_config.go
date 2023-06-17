package configuration

type ConfigApp struct {
	ListenPort         string `split_words:"true" default:":9800"`
	AppName            string `split_words:"true" default:"Smart Contract Service"`
	Version            string `split_words:"true" default:"0.0.1"`
	RootURL            string `split_words:"true" default:"/service/smart-contract"`
	Timeout            int    `split_words:"true" default:"4000"`
	Env                string `split_words:"true" default:"dev"`
	PostgreConnection  string `split_words:"true" default:"host=127.0.0.1 port=5432 dbname=postgres user=postgres password=Sandiaman123. sslmode=disable"`
	SSLMode            string `split_words:"true" default:"disable"`
	LogMode            bool   `split_words:"true" default:"false"`
	RedisConnection    string `split_words:"true" default:"localhost:6379"`
	Secret             string `split_words:"true" default:"rahasia"`
	Expire             int    `split_words:"true" default:"5"`
	RefreshTokenExpire int    `split_words:"true" default:"7"`
	PublicKeyLocation  string `split_words:"true" default:"./assets/rsa256-public.pem"`
	PrivateKeyLocation string `split_words:"true" default:"./assets/rsa256-private.pem"`
}
