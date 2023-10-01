package sqlboiler

import (
	"github.com/volatiletech/sqlboiler/v4/boilingcore"
	"github.com/volatiletech/sqlboiler/v4/drivers"
	"github.com/volatiletech/sqlboiler/v4/importers"
	"os"
)

func NewBoiler() (*boilingcore.State, error) {
	driver, _, err := drivers.RegisterBinaryFromCmdArg("mysql")
	if err != nil {
		return nil, err
	}

	state, err := boilingcore.New(&boilingcore.Config{
		DriverName: driver,
		DriverConfig: map[string]interface{}{
			"blacklist": []string{},
			"whitelist": []string{"users"}, // table name here
			"host":      os.Getenv("DB_HOST"),
			"port":      os.Getenv("DB_PORT"),
			"user":      os.Getenv("DB_USER"),
			"pass":      os.Getenv("DB_PASS"),
			"dbname":    os.Getenv("DB_NAME"),
			"sslmode":   "false",
		},
		Imports:   importers.NewDefaultImports(),
		OutFolder: "/workspace/tmp",
	})
	if err != nil {
		return nil, err
	}

	return state, nil
}
