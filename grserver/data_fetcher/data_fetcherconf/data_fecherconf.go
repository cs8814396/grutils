package data_fetcherconf

type Instance struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Ip       string `toml:"ip"`
	Port     int    `toml:"port"`
}

type Auth struct {
	ClientId     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}

type Conf struct {
	Auth Auth `toml:"auth"`
	//TODO: define you config here
	Instances map[string]Instance `toml:"instances"`
	Querys    map[string]Query    `toml:"querys"`
	Inserts   map[string]Insert   `toml:"inserts"`
}

type Insert struct {
	DatabaseInstance string `toml:"database_instance"`
	DatabaseName     string `toml:"database_name"`
	Statement        string `toml:"statement"`
	//Conditions       map[string]Condition `toml:"conditions"`
	//Delay  bool `toml:"delay"`
	//Upsert bool `toml:"upsert"`
}

type Query struct {
	DatabaseInstance string               `toml:"database_instance"`
	DatabaseName     string               `toml:"database_name"`
	Statement        string               `toml:"statement"`
	Conditions       map[string]Condition `toml:"conditions"`
}

type Condition struct {
	ColumnName      string   `toml:"column_name"`
	PermitOperators []string `toml:"permit_operators"`
	Default         string   `toml:"default"`
}

var GlobalDataFetcherConf Conf
