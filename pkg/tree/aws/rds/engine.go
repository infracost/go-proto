package rds

type RDSEngine uint32

const (
	RDSEngineUnknown           RDSEngine = iota
	RDSEngineMySQL
	RDSEnginePostgres
	RDSEngineMariaDB
	RDSEngineOracleEE
	RDSEngineOracleEECDB
	RDSEngineOracleSE
	RDSEngineOracleSE1
	RDSEngineOracleSE2
	RDSEngineOracleSE2CDB
	RDSEngineSQLServerEE
	RDSEngineSQLServerSE
	RDSEngineSQLServerEX
	RDSEngineSQLServerWeb
	RDSEngineAurora
	RDSEngineAuroraMySQL
	RDSEngineAuroraPostgresql
)

type RDSLicenseModel uint32

const (
	RDSLicenseModelUnknown              RDSLicenseModel = iota
	RDSLicenseModelLicenseIncluded
	RDSLicenseModelBringYourOwnLicense
	RDSLicenseModelGeneralPublicLicense
)
