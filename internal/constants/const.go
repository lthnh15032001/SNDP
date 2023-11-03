package constants

const (
	OffStatus int = iota
	OnStatus
	PROVIDER_AWS                = "aws"
	PROVIDER_GCP                = "gcp"
	PROVIDER_AZURE              = "azure"
	Default                     = ""
	ENV_GORYA_QUEUE_NAME        = "GORYA_QUEUE_NAME"
	ENV_GORYA_REDIS_ADDR        = "GORYA_REDIS_ADDR"
	ENV_GORYA_NUM_WORKER        = "GORYA_NUM_WORKER"
	ENV_GORYA_API_PORT          = "PORT"
	ENV_GORYA_API_HOST          = "HOST"
	ENV_GORYA_ENABLED_PROVIDERS = "GORYA_ENABLED_PROVIDERS"
	ENV_AWS_REGION              = "AWS_REGION"
	ENV_AWS_ENDPOINT            = "AWS_ENDPOINT"
	ENV_GORYA_OIDC_ISSUER_URL   = "GORYA_OIDC_ISSUER_URL"
	ENV_GORYA_OIDC_CLIENT_ID    = "GORYA_OIDC_CLIENT_ID"
	ENV_LOG_LEVEL               = "LOG_LEVEL"
	ENV_GORYA_DB_TYPE           = "GORYA_DB_TYPE"
	ENV_GORYA_DB_HOST           = "GORYA_DB_HOST"
	ENV_GORYA_DB_USER           = "GORYA_DB_USER"
	ENV_GORYA_DB_PASSWORD       = "GORYA_DB_PASSWORD"
	ENV_GORYA_DB_NAME           = "GORYA_DB_NAME"
)
