package sqlite

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	DB *sql.DB
}

func Connect(name string) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", "./storage/"+name+".db")
	if err != nil {
		return nil, err
	}

	// Create table if it doesn't exist
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS meter_usages (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        usage REAL NOT NULL
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &Sqlite{DB: db}, nil
}

// meterusage.csv
// time,usage
func Seed(db *Sqlite) error {
	usages, err := seedCSVData()
	if err != nil {
		return err
	}

	_, err = db.DB.Exec("DELETE FROM meter_usages")
	if err != nil {
		return err
	}

	for _, usage := range usages {
		_, err = db.DB.Exec("INSERT INTO meter_usages (time, usage) VALUES (?, ?)", usage.Time, usage.Usage)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedCSVData() ([]Usage, error) {
	csvFile, err := os.Open("storage/meterusage.csv")
	if err != nil {
		return []Usage{}, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ','

	// skip header
	_, err = reader.Read()
	if err != nil {
		return []Usage{}, err
	}

	var usages []Usage

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []Usage{}, err
		}

		time, err := time.Parse("2006-01-02 15:04:05", record[0])
		if err != nil {
			return []Usage{}, fmt.Errorf("failed to parse time for record %v: %v", record[0], err)
		}

		var usage float64
		if record[1] == "" || record[1] == "NaN" {
			usage = 0.0
		} else {
			usage, err = strconv.ParseFloat(record[1], 64)
			if err != nil {
				return []Usage{}, fmt.Errorf("failed to parse usage: %v", err)
			}
		}

		usages = append(usages, Usage{
			Time:  time,
			Usage: usage,
		})
	}

	return usages, nil
}

func (db *Sqlite) GetUsageHistory(ctx context.Context, startTime, endTime time.Time) ([]Usage, error) {
	rows, err := db.DB.Query("SELECT * FROM meter_usages WHERE time >= ? AND time <= ? ORDER BY time DESC", startTime, endTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	usages := make([]Usage, 0)

	for rows.Next() {
		var usage Usage
		err := rows.Scan(&usage.ID, &usage.Time, &usage.Usage)
		if err != nil {
			return nil, err
		}
		usages = append(usages, usage)
	}

	return usages, nil
}
