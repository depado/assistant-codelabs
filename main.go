package main

import (
	flag "github.com/ogier/pflag"
	"github.com/sirupsen/logrus"

	"github.com/Depado/assistant/backend/conf"
	"github.com/Depado/assistant/backend/database"
	"github.com/Depado/assistant/backend/router"
	"github.com/Depado/assistant/backend/utils"
)

func main() {
	var err error
	var csvf, cnff string
	var expf bool

	flag.StringVarP(&csvf, "import", "i", "", "CSV file to import in the database")
	flag.BoolVarP(&expf, "export", "e", false, "Exports brands, models and such")
	flag.StringVarP(&cnff, "conf", "c", "conf.yml", "Configuration file to import")
	flag.Parse()

	// Load configuration file
	if err = conf.Load(cnff); err != nil {
		logrus.WithError(err).Fatal("Couldn't load configuration file")
	}
	logrus.WithFields(logrus.Fields{"conf": cnff, "listen": conf.C.ListenAddress()}).Info("Parsed configuration file")

	// Initialize the database
	if err = database.Init(); err != nil {
		logrus.WithError(err).Fatal("Couldn't intialize database")
	}

	// Export action
	if expf {
		if err = utils.ExportAll(); err != nil {
			logrus.WithError(err).Fatal("Couldn't export")
		}
		return
	}

	// Import if needed
	if csvf != "" {
		if err = utils.ImportFromFile(csvf); err != nil {
			logrus.WithError(err).Fatal("Couldn't import file")
		}
	}

	// Run the router and start accepting incoming requests
	router.Run()
}
