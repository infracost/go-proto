package database

type Database struct {
	MariaDBServers              []MariaDBServer              `tree:"mariadb_servers"`
	SQLDatabases                []SQLDatabase                `tree:"sql_databases"`
	SQLElasticPools             []SQLElasticPool             `tree:"sql_elastic_pools"`
	SQLManagedInstances         []SQLManagedInstance         `tree:"sql_managed_instances"`
	PostgreSQLServers           []PostgreSQLServer           `tree:"postgresql_servers"`
	PostgreSQLFlexibleServers   []PostgreSQLFlexibleServer   `tree:"postgresql_flexible_servers"`
	MySQLServers                []MySQLServer                `tree:"mysql_servers"`
	MySQLFlexibleServers        []MySQLFlexibleServer        `tree:"mysql_flexible_servers"`
}
