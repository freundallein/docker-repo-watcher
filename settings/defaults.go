package settings

//AppImageName - application docker image name
const AppImageName = "freundallein/drwatcher"

// DefaultRegistryIP - docker registry IP
const DefaultRegistryIP = "192.168.20.126"

// DefaultRegistryPort - docker registry port
const DefaultRegistryPort = "5000"

// DefaultCrontab - crontab for scheduling jobs (first element - seconds)
const DefaultCrontab = "0 0 0 * * *"

// DefaultLogLevel - log level (DEBUG/ERROR)
const DefaultLogLevel = "ERROR"

// DefaultPeriod - image-clean-job period
const DefaultPeriod = 60

// DefaultImageAmount - image-clean-job amount to store
const DefaultImageAmount = 5

// DefaultAutoUpdate - allow drwatcher to autoupdate itself
const DefaultAutoUpdate = true

// DefautlCleanRegistry - allow drwatcher to clean registry
const DefautlCleanRegistry = false
