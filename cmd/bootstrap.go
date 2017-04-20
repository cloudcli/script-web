package cmd

import (
	"script-web/modules/db"
	"script-web/modules/db/models"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	//"script-web/modules/util"
)

var BootstrapCommand = cli.Command{
	Name:        "bootstrap",
	Usage:       "Bootstrap app",
	Description: "Create Tables, Tables seeding",
	Action:      bootstrapAction,
	Subcommands: []cli.Command{
		cli.Command{
			Name:   "create-tables",
			Usage:  "create tables",
			Action: bootstrapCreateTablesCommand,
		},
		cli.Command{
			Name:   "seed",
			Usage:  "create seeding records",
			Action: bootstrapSeedCommand,
		},
	},
}

func bootstrapAction(ctx *cli.Context) error {
	return nil
}

func bootstrapCreateTablesCommand(ctx *cli.Context) error {
	err := createTables()
	return err
}

func bootstrapSeedCommand(ctx *cli.Context) error {
	//adminPass, err := util.EncodePassword("cloudcli_admin")
	//logrus.Info("seeding data records...")
	//if err != nil {
	//	logrus.Errorf("%v", err)
	//	return err
	//}

	//admin, err  := models.NewUser()
	//if err != nil {
	//	logrus.Errorf("%v", err)
	//	return err
	//}
	//admin.Email =
	//admin.Name = "Admin"
	//
	//hashedPassword, err := admin.EncryptPassword("cloudcli_admin")
	//
	//admin.Password = string(hashedPassword)

	mysqldb, err := db.GetDb()
	if err != nil {
		return err
	}

	user, err := mysqldb.CreateUser("admin@cloudcli.com", "Admin", "cloudcli_admin")
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	logrus.Info("create admin user done.", user)
	return nil
}

func createTables() error {
	var tables = []interface{}{
		&models.User{},
		&models.Script{},
		&models.ScriptVersion{},
		&models.Apply{},
		&models.Comment{},
		&models.Message{},
		&models.EventLog{},
		&models.KvStore{},
	}

	mysqldb, err := db.GetDb()
	if err != nil {
		return err
	}
	mysqldb.Db.DropTableIfExists(tables...)
	mysqldb.Db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(tables...)
	return nil
}
