# assistant-codelabs

[![Go Report Card](https://goreportcard.com/badge/github.com/Depado/assistant-codelabs)](https://goreportcard.com/report/github.com/Depado/assistant-codelabs)
![Go Version](https://img.shields.io/badge/go-1.8-brightgreen.svg)
![Go Version](https://img.shields.io/badge/go-1.9-brightgreen.svg)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/bfchroma/blob/master/LICENSE)

This code was written for the Google Assistant Codelabs. It integrates with
Dialogflow.

## Import and Export

The `--import`/`-i` flag with a CSV file allows to import a car dataset, and 
feed it to the main bolt database in a single transaction. 

The CSV format must be in the form :
`id;insertion_date;brand;model;gearbox;fuel;mileage;regdate;price;zip_code;city_name;dpt_name;region_name;type;title;body_md5_hash;word_count`

You can export brands and models in a Dialogflow entity compatible format using
the `--export`/`-e` flag. This will create an `exports/` folder containing
two files : `brands.json` and `models.json`. All you have to do is to copy
the content of these files, switch your Dialogflow entity to raw mode and paste
it.
