package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type CBNRate struct {
	Date     string
	currency string
	buyRate  float64
}

func ParseCSV(reader io.Reader) (*CBNRate, error) {
	r := csv.NewReader(reader)

	r.FieldsPerRecord = -1
	// read header
	_, err := r.Read()

	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV header: %s", err.Error())
	}

	// read latest record (first record)
	rec, err := r.Read()

	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV body: %s", err.Error())
	}

	date := rec[0]
	currency := rec[1]
	buying_rate := rec[4]
	current_date := rec[0]

	for date == current_date {
		if strings.Contains(currency, "US DOLLAR") {
			break
		}
		rec, err = r.Read()
		if err != nil {
			return nil, fmt.Errorf("failed to parse CSV body: %s", err.Error())
		}
		current_date = rec[0]
		currency = rec[1]
		buying_rate = rec[4]
	}

	if strings.Contains(currency, "US DOLLAR") {
		floatVal, err := strconv.ParseFloat(buying_rate, 32)
		if err != nil {
			return nil, err
		}
		return &CBNRate{
			Date:     date,
			currency: "USD",
			buyRate:  floatVal,
		}, nil
	}

	return nil, errors.New("failed to parse cbn rate")
}

func FetchCBNRate() (*CBNRate, error) {
	// fetch from url https://www.cbn.gov.ng/Functions/export.asp?tablename=exchange
	resp, err := http.Get("https://www.cbn.gov.ng/Functions/export.asp?tablename=exchange")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ParseCSV(resp.Body)
}
