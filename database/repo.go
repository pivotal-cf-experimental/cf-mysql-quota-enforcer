package database

import (
	"fmt"

	"database/sql"
	"github.com/pivotal-golang/lager"
)

type Repo interface {
	All() ([]Database, error)
}

type repo struct {
	query  string
	db     *sql.DB
	logger lager.Logger
	logTag string
}

func newRepo(query string, db *sql.DB, logger lager.Logger, logTag string) Repo {
	return &repo{
		query:  query,
		db:     db,
		logger: logger,
		logTag: logTag,
	}
}

func (r repo) All() ([]Database, error) {
	r.logger.Debug(fmt.Sprintf("Executing '%s' database query", r.logTag))

	databases := []Database{}

	rows, err := r.db.Query(r.query)
	if err != nil {
		return databases, fmt.Errorf("Executing '%s' database query: %s", r.logTag, err.Error())
	}
	//TODO: untested Close, due to limitation of sqlmock: https://github.com/DATA-DOG/go-sqlmock/issues/15
	defer rows.Close()

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			//TODO: untested error case, due to limitation of sqlmock: https://github.com/DATA-DOG/go-sqlmock/issues/13
			return databases, fmt.Errorf("Scanning result row of '%s' database query: %s", r.logTag, err.Error())
		}

		databases = append(databases, New(dbName, r.db, r.logger))
	}
	//TODO: untested error case, due to limitation of sqlmock: https://github.com/DATA-DOG/go-sqlmock/issues/13
	if err := rows.Err(); err != nil {
		return databases, fmt.Errorf("Reading result row of '%s' database query: %s", r.logTag, err.Error())
	}

	return databases, nil
}
