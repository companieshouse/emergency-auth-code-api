package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/models"

	_ "github.com/godror/godror"
)

const QUERY = "select od.officer_id, od.officer_forename_1, od.officer_forename_2, od.officer_surname " +
	"from corporate_body cb, corporate_body_appointment cba, officer o, officer_detail od " +
	"where cb.incorporation_number = '%s' " +
	"and cb.corporate_body_id = cba.corporate_body_id " +
	"and cba.officer_id = od.officer_id " +
	"and cba.officer_id = o.officer_id " +
	"and cba.appointment_type_id = 2 " +
	"and cba.resignation_ind = 'N' " +
	"and od.officer_disqualification_ind = 'N' " +
	"and o.corporate_officer_ind = 'N'"

func getOracleDatabase(cfg config.Config) *sql.DB {

	db, err := sql.Open("godror", cfg.OfficerDatabaseUsername+"/"+cfg.OfficerDatabasePassword+"@//"+cfg.OfficerDatabaseUrl)

	// Assume the caller of this func cannot handle the case where there is no database connection
	// so the service must crash here as it cannot continue.
	if err != nil {
		log.Error(err)
	}

	// Check we can connect to oracle database, failure here should crash the program
	if err = db.Ping(); err != nil {
		log.Error(errors.New("ping to oracle timed out. please check the connection to oracle and that it is running: " +
			cfg.OfficerDatabaseUsername + "/" + cfg.OfficerDatabasePassword + "@//" + cfg.OfficerDatabaseUrl))
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("connected to oracle database successfully: %s/%s@//%s", cfg.OfficerDatabaseUsername, cfg.OfficerDatabasePassword, cfg.OfficerDatabaseUrl))

	return db
}

// OracleService is an implementation of the OfficerDAOService interface
type OracleService struct {
	db *sql.DB
}

// GetCompanyOfficers returns a list of valid company officers
func (o OracleService) GetCompanyOfficers(companyNumber string) (*models.CompanyOfficers, error) {

	// run query for officers on connected database
	rows, err := o.db.Query(fmt.Sprintf(QUERY, companyNumber))
	if err != nil {
		return nil, errors.New("failed to run query for company officers on database")
	}

	defer rows.Close()

	// parse through oracle response and transform to JSON to be returned
	companyOfficers, err := oracleCompanyOfficersToCompanyOfficers(rows)
	if err != nil {
		return nil, err
	}

	return companyOfficers, nil
}

// oracleCompanyOfficersToCompanyOfficers transforms incoming oracle response for company officers to readable JSON
func oracleCompanyOfficersToCompanyOfficers(rows *sql.Rows) (*models.CompanyOfficers, error) {

	companyOfficers := &models.CompanyOfficers{}

	for rows.Next() {
		officer := models.Items{}
		if err := rows.Scan(&officer.ID, &officer.Forename1, &officer.Forename2, &officer.Surname); err != nil {
			return nil, errors.New("error reading rows returned from database query")
		}

		companyOfficers.Items = append(companyOfficers.Items, officer)
	}

	companyOfficers.TotalCount = len(companyOfficers.Items)

	return companyOfficers, nil
}
