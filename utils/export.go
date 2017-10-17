package utils

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Depado/assistant/backend/database"
	"github.com/Depado/assistant/backend/dialogflow"
	"github.com/Depado/assistant/backend/models"
)

// ExportAll will export some fields to JSON files in the dialogflow entity
// format
func ExportAll() error {
	var err error
	var cars []models.Car

	s := time.Now()
	logrus.WithField("types", []string{"brands", "models"}).Info("Export started")

	if err = database.DB.All(&cars); err != nil {
		return errors.Wrap(err, "ExportBrands(): couldn't query all cars")
	}

	brands := make(map[string]bool)
	for _, c := range cars {
		brands[c.Brand] = true
	}
	be := dialogflow.Entities{}
	for k := range brands {
		if k != "" {
			be = append(be, dialogflow.Entity{Value: k, Synonyms: []string{k}})
		}
	}
	out, err := json.Marshal(&be)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile("exports/brands.json", out, 0644); err != nil {
		return err
	}

	mod := make(map[string]bool)
	for _, c := range cars {
		mod[c.Model] = true
	}
	me := dialogflow.Entities{}
	for k := range mod {
		if k != "" {
			me = append(me, dialogflow.Entity{Value: k, Synonyms: []string{k}})
		}
	}
	out, err = json.Marshal(&me)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile("exports/models.json", out, 0644); err != nil {
		return err
	}

	logrus.WithField("took", time.Since(s)).Info("Done exporting")
	return nil
}
