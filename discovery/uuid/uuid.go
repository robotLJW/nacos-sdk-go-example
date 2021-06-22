package uuid

import "github.com/google/uuid"

func GenerateServiceName() string {
	serviceName, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return serviceName.String()
}
