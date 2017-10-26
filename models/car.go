package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Depado/assistant-codelabs/database"
	"github.com/Depado/assistant-codelabs/dialogflow"
)

// Car is a simple struct holding car information
type Car struct {
	StoreID       int64     `json:"store_id,omitempty" storm:"id"`
	InsertionDate time.Time `json:"insertion_date,omitempty"`
	Brand         string    `json:"brand,omitempty" storm:"index"`
	Model         string    `json:"model,omitempty" storm:"index"`
	Gearbox       string    `json:"gearbox,omitempty" storm:"index"`
	Fuel          string    `json:"fuel,omitempty" storm:"index"`
	Mileage       int       `json:"mileage,omitempty"`
	RegDate       int       `json:"reg_date,omitempty"`
	Price         float64   `json:"price,omitempty"`
	ZipCode       int       `json:"zip_code,omitempty"`
	City          string    `json:"city,omitempty"`
	Department    string    `json:"department,omitempty"`
	Region        string    `json:"region,omitempty"`
	UploadType    string    `json:"upload_type,omitempty"`
	Title         string    `json:"title,omitempty"`
	BodyMD5       string    `json:"body_md_5,omitempty"`
	WordCount     int       `json:"word_count,omitempty"`
	MaxListID     int64     `json:"max_list_id,omitempty"`
}

// Save saves the given car to the database
func (c Car) Save() error {
	return database.DB.Save(&c)
}

// NewFromCSV parses a single CSV line into an actual Car struct
func NewFromCSV(line []string) (*Car, error) {
	var err error
	var id, mileage, reg, zipcode, wordcount, maxlid int
	var price float64
	var t time.Time

	if id, err = strconv.Atoi(line[0]); err != nil {
		return &Car{}, errors.Wrap(err, "NewFromCSV(): unable to parse store_id")
	}
	if mileage, err = strconv.Atoi(line[6]); err != nil {
		return &Car{}, errors.Wrap(err, "NewFromCSV(): unable to parse mileage")
	}
	if reg, err = strconv.Atoi(line[7]); err != nil {
		return &Car{}, errors.Wrap(err, "NewFromCSV(): unable to parse reg_date")
	}
	if price, err = strconv.ParseFloat(line[8], 64); err != nil {
		return &Car{}, errors.Wrap(err, "NewFromCSV(): unable to parse price")
	}
	if zipcode, err = strconv.Atoi(line[9]); err != nil {
		return &Car{}, errors.Wrap(err, "NewFromCSV(): unable to parse reg_date")
	}
	if wordcount, err = strconv.Atoi(line[16]); err != nil {
		return &Car{}, errors.Wrap(err, "NewFromCSV(): unable to parse word_count")
	}
	if maxlid, err = strconv.Atoi(line[17]); err != nil {
		return &Car{}, errors.Wrap(err, "NewFromCSV(): unable to parse max_list_id")
	}
	if t, err = time.Parse("2006-01-02", line[1]); err != nil {
		return &Car{}, errors.Wrap(err, "NewFromCSV(): unable to parse insertion_date")
	}
	return &Car{
		StoreID:       int64(id),
		InsertionDate: t,
		Brand:         line[2],
		Model:         line[3],
		Gearbox:       line[4],
		Fuel:          line[5],
		Mileage:       mileage,
		RegDate:       reg,
		Price:         price,
		ZipCode:       zipcode,
		City:          line[10],
		Department:    line[11],
		Region:        line[12],
		UploadType:    line[13],
		Title:         line[14],
		BodyMD5:       line[15],
		WordCount:     wordcount,
		MaxListID:     int64(maxlid),
	}, nil
}

// NewCarFromParameters returns a new car from dialogflow's parameters
func NewCarFromParameters(params dialogflow.CarSearchParameters) (*Car, error) {

	c := &Car{
		Brand:   params.CarBrand,
		Model:   params.CarModel,
		Fuel:    params.CarEnergy,
		Gearbox: params.CarGearbox,
		Mileage: params.CarKilometers.Amount,
	}
	c.RegDate, _ = strconv.Atoi(params.Year)

	return c, nil
}

// GetEstimate estimates the price of a car
func (c Car) GetEstimate() (float64, error) {
	var err error
	var cars []Car

	queries := []q.Matcher{}

	if c.Brand != "" {
		queries = append(queries, q.Eq("Brand", c.Brand))
	}
	if c.Model != "" {
		queries = append(queries, q.Eq("Model", c.Model))
	}
	if c.Fuel != "" {
		queries = append(queries, q.Or(q.Eq("Fuel", c.Fuel), q.Eq("Fuel", strings.Title(c.Fuel))))
	}
	if c.Gearbox != "" {
		queries = append(queries, q.Or(q.Eq("Gearbox", c.Gearbox), q.Eq("Gearbox", strings.Title(c.Gearbox))))
	}
	if c.Mileage != 0 {
		queries = append(queries, q.Gt("Mileage", c.Mileage-10000), q.Lt("Mileage", c.Mileage+10000))
	}
	if c.RegDate != 0 {
		queries = append(queries, q.Gt("RegDate", c.RegDate-1), q.Lt("RegDate", c.RegDate+1))
	}
	if err = database.DB.Select(queries...).Find(&cars); err != nil {
		if err == storm.ErrNotFound {
			return 0, errors.Wrap(err, "GetEstimate(): no entry found")
		}
	}
	var tot float64
	for _, cc := range cars {
		tot += cc.Price
	}
	estimate := tot / float64(len(cars))
	logrus.WithFields(logrus.Fields{"matched": len(cars), "average": estimate}).Info("Done querying")
	return estimate, nil
}

// Summarize returns a string representing the car
func (c Car) Summarize() string {
	s := ""
	if c.Brand != "" {
		s = c.Brand
	}
	if c.Model != "" {
		if len(s) == 0 {
			s = c.Model
		} else {
			s = fmt.Sprintf("%s %s", s, c.Model)
		}
	}
	if c.Fuel != "" {
		if len(s) == 0 {
			s = fmt.Sprintf("voiture %s", c.Fuel)
		} else {
			s = fmt.Sprintf("%s %s", s, c.Fuel)
		}
	}
	if c.RegDate != 0 {
		if len(s) == 0 {
			s = fmt.Sprintf("voiture de %d", c.RegDate)
		} else {
			s = fmt.Sprintf("%s de %d", s, c.RegDate)
		}
	}
	if c.Gearbox != "" {
		if len(s) == 0 {
			s = fmt.Sprintf("une voiture avec une boite %s", strings.ToLower(c.Gearbox))
		} else {
			s = fmt.Sprintf("%s en boite %s", s, strings.ToLower(c.Gearbox))
		}
	}
	if c.Mileage != 0 {
		if len(s) == 0 {
			s = fmt.Sprintf("une voiture avec %d kilomètres au compteur", c.Mileage)
		} else {
			s = fmt.Sprintf("%s avec %d kilomètres au compteur", s, c.Mileage)
		}
	}
	return s
}

// Estimate gets the estimated value and formats it
func (c Car) Estimate() string {
	var err error
	var p float64

	if p, err = c.GetEstimate(); err != nil {
		return fmt.Sprintf("Je n'ai pas trouvé de prix moyen pour une %s", c.Summarize())
	}
	return fmt.Sprintf("Une %s se vend environ %d€ sur leboncoin", c.Summarize(), int(p))
}
