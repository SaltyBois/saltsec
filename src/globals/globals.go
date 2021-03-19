package globals

import (
	"os"
	"strconv"
)

// TODO(Jovan): Use envvar for port
const (
	DEFAULT_PORT        = ":8081"
	DEFAULT_HOST_DB     = "localhost"
	DEFAULT_PORT_DB     = 5432
	DEFAULT_USER_DB     = "postgres"
	DEFAULT_PASSWORD_DB = "root"
	DEFAULT_NAME_DB     = "SaltyData"
)

var (
	PORT        string
	HOST_DB     string
	PORT_DB     int
	USER_DB     string
	PASSWORD_DB string
	NAME_DB     string
)

// TODO(Jovan): Make it cleaner???
func LoadGlobalVars() {
	if val, present := os.LookupEnv("PORT_SALT"); present {
		PORT = val
	} else {
		PORT = DEFAULT_PORT
	}

	if val, present := os.LookupEnv("HOST_DB_SALT"); present {
		HOST_DB = val
	} else {
		HOST_DB = DEFAULT_HOST_DB
	}

	if val, present := os.LookupEnv("PORT_DB_SALT"); present {
		PORT_DB, _ = strconv.Atoi(val)
	} else {
		PORT_DB = DEFAULT_PORT_DB
	}

	if val, present := os.LookupEnv("USER_DB_SALT"); present {
		USER_DB = val
	} else {
		USER_DB = DEFAULT_USER_DB
	}

	if val, present := os.LookupEnv("PASSWORD_DB_SALT"); present {
		PASSWORD_DB = val
	} else {
		PASSWORD_DB = DEFAULT_PASSWORD_DB
	}

	if val, present := os.LookupEnv("NAME_DB_SALT"); present {
		NAME_DB = val
	} else {
		NAME_DB = DEFAULT_NAME_DB
	}
}
