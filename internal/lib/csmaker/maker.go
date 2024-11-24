package csmaker

import "fmt"

func MakeConnectionString(user, password, dbName string, options ...map[string]string) string {
	host := "localhost"
	port := "5432"
	sslmode := "disable"
	sslrootcert := ""

	if len(options) > 0 {
		for _, option := range options {
			if val, ok := option["host"]; ok && val != "" {
				host = val
			}
			if val, ok := option["port"]; ok && val != "" {
				port = val
			}
			if val, ok := option["sslmode"]; ok && val != "" {
				sslmode = val
			}
			if val, ok := option["sslrootcert"]; ok && val != "" {
				sslrootcert = fmt.Sprintf("sslrootcert=%s", val)
			}
		}
	}

	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s %s",
		user, password, host, port, dbName, sslmode, sslrootcert)
	return connectionString
}
