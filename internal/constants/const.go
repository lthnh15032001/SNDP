package constants

const (
	OffStatus int = iota
	OnStatus
	PROVIDER_AWS          = "aws"
	PROVIDER_GCP          = "gcp"
	PROVIDER_AZURE        = "azure"
	Default               = ""
	ENV_QUEUE_NAME        = "QUEUE_NAME"
	ENV_REDIS_ADDR        = "REDIS_ADDR"
	ENV_NUM_WORKER        = "NUM_WORKER"
	ENV_API_PORT          = "PORT"
	ENV_API_HOST          = "HOST"
	ENV_ENABLED_PROVIDERS = "ENABLED_PROVIDERS"
	ENV_AWS_REGION        = "AWS_REGION"
	ENV_AWS_ENDPOINT      = "AWS_ENDPOINT"
	ENV_OIDC_ISSUER_URL   = "OIDC_ISSUER_URL"
	ENV_OIDC_CLIENT_ID    = "OIDC_CLIENT_ID"
	ENV_LOG_LEVEL         = "LOG_LEVEL"
	ENV_DB_TYPE           = "db.type"
	ENV_DB_HOST           = "db.host"
	ENV_DB_USER           = "db.user"
	ENV_DB_PASSWORD       = "db.password"
	ENV_DB_NAME           = "db.name"
)
