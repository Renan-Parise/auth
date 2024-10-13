package utils

import "os"

func GetMailServiceURL() string {
	return os.Getenv("MAIL_SERVICE_URL")
}
