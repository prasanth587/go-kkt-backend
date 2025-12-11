package mssqlcon

import (
	"context"
	"database/sql"
	"fmt"
	lg "log"
	"sync"

	"github.com/prabha303-vi/log-util/log"
	//_ "github.com/hgfischer/mysql"

	"go-transport-hub/constant"
	"go-transport-hub/dtos/schema"
	"go-transport-hub/utils"

	_ "github.com/go-sql-driver/mysql"
)

// Pool of database connection
var once sync.Once
var ConnPool *sql.DB

type DBConn struct {
	conn          *sql.DB
	tx            *sql.Tx
	isTransaction bool
	l             *log.Logger
}
type IDBConn interface {
	Init(*log.Logger)
	GetQueryer() Queryer
	ExecuteInTransaction(func() error) error
	// rollbackTransaction(tx *sql.Tx)
}

// Interface to abstract the queryer(dbconnection or transaction)
type Queryer interface {
	Exec(sql string, arguments ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
	Query(sql string, args ...interface{}) (*sql.Rows, error)
	QueryRow(sql string, args ...interface{}) *sql.Row
	Prepare(sql string) (*sql.Stmt, error)
}

func New(l *log.Logger) *DBConn {
	return &DBConn{
		l:    l,
		conn: ConnPool,
	}
}

func NewDBConn(l *log.Logger, conn *sql.DB) *DBConn {
	return &DBConn{
		l:    l,
		conn: conn,
	}
}

// Initialize the DB connection and assign the existing db connection
func (db *DBConn) Init(l *log.Logger) {
	db.conn = ConnPool
	db.l = l
}

func (db *DBConn) GetQueryer() Queryer {
	if db.isTransaction {
		return db.tx
	} else {
		return db.conn
	}
}

// ExecuteInTransaction executes the given function in DB transaction, i.e. It commits
// only if there is not error otherwise it is rolledback.
func (db *DBConn) ExecuteInTransaction(f func() error) (err error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	db.tx = tx
	db.isTransaction = true

	defer func() {
		if r := recover(); r != nil {
			db.l.Fatal("Recovered in function ", r)
			db.rollbackTransaction(tx)
		}
		db.isTransaction = false
	}()

	err = f()
	if err != nil {
		db.rollbackTransaction(tx)
		return err
	}
	err = tx.Commit()
	if err != nil {
		db.rollbackTransaction(tx)
		return err
	}
	return nil
}

func (db *DBConn) rollbackTransaction(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		db.l.Error("Error While rollback, Err: ", err)
	}
}

func MSSqlInit(url string) {
	if ConnPool == nil {
		once.Do(func() {
			var err error
			// Create connection pool
			ConnPool, err = sql.Open("mysql", url)
			if err != nil {
				lg.Printf("Error creating connection pool: %+v \n", err)
			}
			ConnPool.SetMaxOpenConns(15)
			ctx := context.Background()
			err = ConnPool.PingContext(ctx)
			if err != nil {
				lg.Printf("Unable to ping to DB. Err: %+v", err)
				return
			}
			lg.Println("Connected to database successfully!")
			CreateTable(ConnPool)
			CreateDefaultValues(ConnPool)
		})
	}
}

func CreateTable(db *sql.DB) {
	createOrganisation(db)
	createRole(db)
	createRolePermission(db)
	createWebScreens(db)
	createRoleAndScreensMap(db)
	createUser(db)
	createEmployee(db)
	createDriver(db)
	createCustomerInvoice(db)
	createVehicleSizeTypes(db)
	//createVehicle(db)
	//createVendor(db)
	createVendors(db)
	createVehicles(db)
	createTruckTypes(db)
	createMaterialTypes(db)
	createBranch(db)
	//createCustomer(db)
	createCustomers(db)
	createLoadingPoint(db)
	createtripSheetNumer(db)
	createTripSheetHeader(db)
	createTripSheetHeaderLoadingPoints(db)
	createManagePod(db)
	createLRHeader(db)

	createCustomersContactInfo(db)
	createCustomersAggrement(db)

	createDriverContactInfo(db)
	createDeclarationYearInfo(db)
	createAppConfig(db)
	createEmployeeAttendance(db)
	createEmployeeAttendanceEntry(db)
	createNotifications(db)
}

func createNotifications(db *sql.DB) {
	notificationTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		notifications (
			notification_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			org_id INT UNSIGNED NOT NULL,
			CONSTRAINT fky_notification_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
			user_id INT UNSIGNED,
			CONSTRAINT fky_notification_user FOREIGN KEY (user_id) REFERENCES user_login(id),
			notification_type VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			related_entity_type VARCHAR(50),
			related_entity_id INT UNSIGNED,
			is_read BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (notification_id),
			INDEX idx_org_user (org_id, user_id),
			INDEX idx_is_read (is_read),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println("ERROR: notifications table prepare: ", err.Error())
		return
	}
	_, err = notificationTable.Exec()
	if err != nil {
		lg.Println("ERROR: notifications table: ", err.Error())
	}
}

func CreateDefaultValues(db *sql.DB) {
	insertSuperAdmin(db)
	insertTruckTypes(db)
	insertMaterialTypes(db)
	insertRolePermisstionLabel(db)
	insertwebsiteScreens(db)
	insertAppConfig(db)
}

func MSSqlConnClose() {
	if ConnPool != nil {
		ConnPool.Close()
	}
}

func createOrganisation(db *sql.DB) {

	organisationTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		organisation (org_id int unsigned NOT NULL AUTO_INCREMENT, 
		name varchar(50) NOT NULL, 
		display_name varchar(50), 
		domain_name varchar(50),
		email_id varchar(50),
		contact_name varchar(50),
		contact_no varchar(50),
		is_active BOOLEAN DEFAULT TRUE,
		logo_path varchar(500),
		address_line1 TEXT,
		address_line2 TEXT,
		city varchar(100),
		state varchar(100),
		country varchar(100),
		zipcode varchar(6),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		version INT,
		CONSTRAINT org_constrain UNIQUE (email_id),
		PRIMARY KEY (org_id))ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
		return
	}
	_, err = organisationTable.Exec()
	if err != nil {
		lg.Println("ERROR: organisation table: ", err.Error())
	}

}
func createRole(db *sql.DB) {

	roleTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		thub_role ( role_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		role_code VARCHAR(50) NOT NULL,
		role_name VARCHAR(100) NOT NULL,
		description VARCHAR(255) NOT NULL,
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_role_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		is_active BOOLEAN DEFAULT TRUE,
		version INT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT codex_role UNIQUE (role_code, org_id),
		PRIMARY KEY (role_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = roleTable.Exec()
	if err != nil {
		lg.Println("ERROR: Role table: ", err.Error())
	}
}

func createRoleAndScreensMap(db *sql.DB) {

	table, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		thub_role_screens_map ( thub_role_screens_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		role_code VARCHAR(50),
		role_name VARCHAR(100),
		role_id INT UNSIGNED,
 		CONSTRAINT fky_role_id_map FOREIGN KEY (role_id) REFERENCES thub_role(role_id),
		website_screen_id INT UNSIGNED,
 		CONSTRAINT fky_website_screen FOREIGN KEY (website_screen_id) REFERENCES website_screens(website_screen_id),
		permisstion_label_id INT UNSIGNED,
 		CONSTRAINT fky_permisstion_label FOREIGN KEY (permisstion_label_id) REFERENCES permisstion_label(permisstion_label_id),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT codex_role_screens UNIQUE (role_id, website_screen_id),
		CONSTRAINT codex_role_screens_permission UNIQUE (role_id, website_screen_id, permisstion_label_id),
		PRIMARY KEY (thub_role_screens_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = table.Exec()
	if err != nil {
		lg.Println("ERROR: thub_role_screens_map table: ", err.Error())
	}

}

func createRolePermission(db *sql.DB) {
	roleTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		permisstion_label ( permisstion_label_id INT UNSIGNED NOT NULL AUTO_INCREMENT,		
		permisstion_label VARCHAR(50) NOT NULL,
		description VARCHAR(255) NOT NULL,
 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,		
		PRIMARY KEY (permisstion_label_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = roleTable.Exec()
	if err != nil {
		lg.Println("ERROR: Role table: ", err.Error())
	}
}

func createWebScreens(db *sql.DB) {
	roleTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		website_screens (website_screen_id INT UNSIGNED NOT NULL AUTO_INCREMENT,		
		screen_name VARCHAR(50) NOT NULL,
		description VARCHAR(255) NOT NULL,
 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,		
		PRIMARY KEY (website_screen_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = roleTable.Exec()
	if err != nil {
		lg.Println("ERROR: Role table: ", err.Error())
	}
}

func createUser(db *sql.DB) {
	createUser, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
	user_login (id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		first_name VARCHAR(25) NOT NULL,
		last_name VARCHAR(25) NOT NULL,
		mobile_no VARCHAR(20),
		email_id VARCHAR(50),
		password VARCHAR(500),
		password_str VARCHAR(50),
		is_active BOOLEAN DEFAULT TRUE,
		access_token VARCHAR(500),
		role_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_user_const FOREIGN KEY (role_id) REFERENCES thub_role(role_id),
		employee_id INT,
		is_super_admin BOOLEAN DEFAULT FALSE,
		is_admin BOOLEAN DEFAULT FALSE,
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_user_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		login_type VARCHAR(10) NOT NULL,
		last_login DATETIME,
		facility_id INT,
		version INT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT userx_con UNIQUE (mobile_no, email_id),
		PRIMARY KEY (id)
		) ENGINE=InnoDB;`)

	if err != nil {
		lg.Println(err.Error())
	}
	_, err = createUser.Exec()
	if err != nil {
		lg.Println("ERROR: user table: ", err.Error())
	}
}

func createEmployee(db *sql.DB) {
	createUser, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
	employee (emp_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		first_name VARCHAR(25) NOT NULL,
		last_name VARCHAR(25),
		employee_code VARCHAR(20) NOT NULL,
		CONSTRAINT empx_code UNIQUE (employee_code),
		mobile_no VARCHAR(20) NOT NULL,
		email_id VARCHAR(50),
		role_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_emp_const FOREIGN KEY (role_id) REFERENCES thub_role(role_id),
		dob VARCHAR(30),
		gender VARCHAR(10),
		aadhar_no VARCHAR(20),
		access_no VARCHAR(20),
		is_active BOOLEAN DEFAULT TRUE,
		access_token VARCHAR(500),
		joining_date VARCHAR(30),
		relieving_date VARCHAR(30),
		employee_performance VARCHAR(50),
		vehicle_assigned VARCHAR(50),
		monthly_salary INT,
		annual_salary INT,
		annual_bonus INT,
		address_line1 TEXT,
		address_line2 TEXT,
		city varchar(100),
		state varchar(100),
		country varchar(100),
		is_super_admin BOOLEAN DEFAULT FALSE,
		is_admin BOOLEAN DEFAULT FALSE,
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_emp_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		login_type VARCHAR(10) NOT NULL,
		image varchar(1000),
		facility_id INT,
		version INT,
		pin_code INT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT empx_con UNIQUE (mobile_no, email_id),
		PRIMARY KEY (emp_id)
		) ENGINE=InnoDB;`)

	if err != nil {
		lg.Println(err.Error())
	}
	_, err = createUser.Exec()
	if err != nil {
		lg.Println("ERROR: empoyee table: ", err.Error())
	}
}

/* func createVendor(db *sql.DB) {
	vendor, err := db.Prepare(`CREATE TABLE IF NOT EXISTS vendor (
		vendor_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		vendor_name VARCHAR(100) NOT NULL,
		vendor_code VARCHAR(50) NOT NULL,
		mobile_number VARCHAR(20) NOT NULL,
		contact_person VARCHAR(50) NOT NULL,
		alternative_number VARCHAR(20),
		status VARCHAR(20) NOT NULL,
		is_active BOOLEAN DEFAULT TRUE NOT NULL,
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_vendor_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		login_type VARCHAR(10),
		address_line1 TEXT,
		address_line2 TEXT,
		city VARCHAR(100),
		state VARCHAR(100),
		visiting_card_image VARCHAR(1000),
		pancard_img VARCHAR(1000),
		aadhar_card_img VARCHAR(1000),
		cancelled_check_book_img VARCHAR(1000),
		bank_passbook_img VARCHAR(1000),
		gst_document_img VARCHAR(1000),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (vendor_id),
		CONSTRAINT vendorx_mobile_exists UNIQUE (mobile_number),
		CONSTRAINT vendorx_vendor_code_exists UNIQUE (vendor_code)
	)ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = vendor.Exec()
	if err != nil {
		lg.Println("ERROR: vendor table: ", err.Error())
	}
} */

/*func createVehicle(db *sql.DB) {
	driver, err := db.Prepare(`CREATE TABLE IF NOT EXISTS
	vehicle (
	vehicle_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    vehicle_type VARCHAR(50) NOT NULL,
    vehicle_number VARCHAR(20) NOT NULL,
    vehicle_model VARCHAR(50) NOT NULL,
    vehicle_year INT NOT NULL,
    vehicle_capacity VARCHAR(200),
    vehicle_insurance_number VARCHAR(50),
    insurance_expiry_date VARCHAR(20),
    vehicle_registration_date VARCHAR(30),
    vehicle_renewal_date VARCHAR(30),
    is_active BOOLEAN DEFAULT TRUE,
	status VARCHAR(20) NOT NULL,
    driver_id INT UNSIGNED,
    org_id INT UNSIGNED NOT NULL,
    vehicle_image VARCHAR(1000),
    fitness_certificate VARCHAR(1000),
    insurance_certificate VARCHAR(1000),
	pollution_certificate VARCHAR(1000),
	national_permits_certificate VARCHAR(1000),
	registration_certificate VARCHAR(1000),
	annual_maintenance_certificate VARCHAR(1000),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_vehicle_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
    CONSTRAINT fk_vehicle_driver FOREIGN KEY (driver_id) REFERENCES driver(driver_id),
    CONSTRAINT vehiclex_vehicle_no_exists UNIQUE (vehicle_number),
    PRIMARY KEY (vehicle_id)
    ) ENGINE=InnoDB;`)

	if err != nil {
		lg.Println(err.Error())
	}
	_, err = driver.Exec()
	if err != nil {
		lg.Println("ERROR: driver table: ", err.Error())
	}
} */

func createDriver(db *sql.DB) {
	driver, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
	driver (driver_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		first_name VARCHAR(25) NOT NULL,
		last_name VARCHAR(25),
		license_number VARCHAR(50) NOT NULL,
		license_expiry_date VARCHAR(20) NOT NULL,
		mobile_no VARCHAR(20) NOT NULL,
		alternate_contact_number VARCHAR(20) NOT NULL,
		email_id VARCHAR(50),
		joining_date VARCHAR(30),
		relieving_date VARCHAR(30),
		is_active BOOLEAN DEFAULT TRUE,
		vehicle_id INT,
		driver_experience DECIMAL(5, 2),
		address_line1 TEXT,
		address_line2 TEXT,
		city varchar(100),
		state varchar(100),
		country varchar(100),
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_driver_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		login_type VARCHAR(10) NOT NULL,
		license_front_img varchar(1000),
		license_back_img varchar(1000),
		other_document varchar(1000),
		image varchar(1000),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		CONSTRAINT driverx_mobile_exists UNIQUE (mobile_no),
		CONSTRAINT driverx_lic_no_exists UNIQUE (license_number),
		PRIMARY KEY (driver_id)
		) ENGINE=InnoDB;`)

	if err != nil {
		lg.Println(err.Error())
	}
	_, err = driver.Exec()
	if err != nil {
		lg.Println("ERROR: driver table: ", err.Error())
	}
}

/* func createCustomer(db *sql.DB) {
	customer, err := db.Prepare(`CREATE TABLE IF NOT EXISTS customer (
    customer_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	branch_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_customer_branch_ref FOREIGN KEY (branch_id) REFERENCES branch(branch_id),
	customer_code VARCHAR(50) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
	address TEXT NOT NULL,
	city VARCHAR(100) NOT NULL,
	state VARCHAR(100),
	pincode INT UNSIGNED NOT NULL,
	country varchar(100),
	gstin_type varchar(100),
	gstin_no varchar(100),
	pan_number varchar(100),
	contact_person VARCHAR(100) NOT NULL,
	mobile_number VARCHAR(20) UNIQUE NOT NULL,
	alternative_number VARCHAR(20),
	email_id VARCHAR(100) UNIQUE NOT NULL,
	status VARCHAR(20),
	branch_sales_person_name varchar(100),
	branch_responsible_person varchar(100),
	employee_code varchar(100),
	agreement_doc_image VARCHAR(1000),
	business_start_date VARCHAR(1000),
	is_active BOOLEAN DEFAULT TRUE,
	org_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_customer_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
    credibility_days INT,
	company_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = customer.Exec()
	if err != nil {
		lg.Println("ERROR: customer table: ", err.Error())
	}
} */

func createTruckTypes(db *sql.DB) {

	ttTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		truck_types (truck_type_id int unsigned NOT NULL AUTO_INCREMENT, 
		category varchar(50) NOT NULL, 
		truck_type varchar(50) NOT NULL, 
		is_active BOOLEAN DEFAULT TRUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		org_id INT UNSIGNED,
		CONSTRAINT fky_tt_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		PRIMARY KEY (truck_type_id))ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
		return
	}
	_, err = ttTable.Exec()
	if err != nil {
		lg.Println("ERROR: truck_types table: ", err.Error())
	}
}

func createMaterialTypes(db *sql.DB) {

	mtTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS material_types (
    material_type_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    material_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    org_id INT UNSIGNED,
    CONSTRAINT fky_mt_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
    PRIMARY KEY (material_type_id)
    ) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
		return
	}
	_, err = mtTable.Exec()
	if err != nil {
		lg.Println("ERROR: material_types table: ", err.Error())
	}
}

func createBranch(db *sql.DB) {
	branchTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		branch (branch_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		branch_name VARCHAR(50) NOT NULL,
		branch_code VARCHAR(100) UNIQUE NOT NULL,
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_branch_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		is_active BOOLEAN DEFAULT TRUE,
		address_line1 TEXT,
		address_line2 TEXT,
		city varchar(100),
		state varchar(100),
		country varchar(100),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (branch_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = branchTable.Exec()
	if err != nil {
		lg.Println("ERROR: branch table: ", err.Error())
	}
}

func createLoadingPoint(db *sql.DB) {
	branchTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		loading_point (loading_point_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		branch_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_loading_branch FOREIGN KEY (branch_id) REFERENCES branch(branch_id),
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_loading_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		is_active BOOLEAN DEFAULT TRUE,
		city_code varchar(50) UNIQUE NOT NULL,
		city_name varchar(50) NOT NULL,
		address_line TEXT,
		map_link varchar(1000) NOT NULL,
		state varchar(100) NOT NULL,
		country varchar(50),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (loading_point_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = branchTable.Exec()
	if err != nil {
		lg.Println("ERROR: branch table: ", err.Error())
	}
}

func createtripSheetNumer(db *sql.DB) {
	tripNumTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		trip_sheet_num (id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		year varchar(50) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = tripNumTable.Exec()
	if err != nil {
		lg.Println("ERROR: branch table: ", err.Error())
	}
}

func createTripSheetHeader(db *sql.DB) {
	tripSheet, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		trip_sheet_header (trip_sheet_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		trip_sheet_num varchar(100) UNIQUE NOT NULL,
		trip_type varchar(50) NOT NULL,
		trip_sheet_type varchar(100) NOT NULL,
		load_hours_type varchar(100) NOT NULL,
		open_trip_date_time varchar(100) NOT NULL,

		branch_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_tripsht_ref FOREIGN KEY (branch_id) REFERENCES branch(branch_id),

		customer_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_customer_ref FOREIGN KEY (customer_id) REFERENCES customers (customer_id),
		
		vendor_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_vendor_ref FOREIGN KEY (vendor_id) REFERENCES vendors (vendor_id),

		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_tripsheet_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),

		vehicle_size_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_vehicle_size_type FOREIGN KEY (vehicle_size_id) REFERENCES vehicle_size_type(vehicle_size_id),

		created_by INT UNSIGNED NOT NULL,
		CONSTRAINT fky_created_by_rel FOREIGN KEY (created_by) REFERENCES user_login(id),

		managed_by INT UNSIGNED NOT NULL,
		CONSTRAINT fky_managed_by_rel FOREIGN KEY (managed_by) REFERENCES user_login(id),

		zonal_name varchar(100),

		 vehicle_capacity_ton varchar(100),
		 vehicle_number varchar(100),
		 vehicle_size varchar(100),
		 mobile_number varchar(50),
		 driver_name varchar(50),
		 driver_license_image varchar(1000),
		 lr_gate_image varchar(1000),
		 lr_number varchar(50) UNIQUE,
		 CONSTRAINT unique_lr_number UNIQUE (lr_number),

		 trip_submitted_date varchar(50),
		 trip_closed_date varchar(50),
		 trip_delivered_date varchar(50),
		 trip_completed_date varchar(50),

		 customer_invoice_no varchar(50),
		 customer_close_trip_date_time varchar(50),
		 customer_payment_received_date varchar(50),
		 customer_remark varchar(50),
		 customer_billing_raised_date varchar(50),

		 customer_base_rate NUMERIC(10, 2),
		 customer_km_cost NUMERIC(10, 2),
		 customer_toll NUMERIC(10, 2),
		 customer_extra_hours  NUMERIC(10, 2),
		 customer_extra_km  NUMERIC(10, 2),
		 customer_total_hire NUMERIC(10, 2),
		 customer_debit_amount NUMERIC(10, 2),

		 customer_per_load_hire NUMERIC(10, 2),
		 customer_running_km NUMERIC(10, 2),
		 customer_per_km_price NUMERIC(10, 2),
		 customer_placed_vehicle_size varchar(100),
		 customer_load_cancelled varchar(100),
		 customer_reported_date_time_for_halting_calc varchar(50),

		 customer_invoice_id BIGINT UNSIGNED,
		 CONSTRAINT fky_customer_invoice_rel FOREIGN KEY (customer_invoice_id) REFERENCES customer_invoice(id),
		  
		 pod_required INT DEFAULT 0,
		 CONSTRAINT tp_pod_required_0_or_1 CHECK (pod_required IN (0, 1)),

		 vendor_paid_date varchar(50),
		 vendor_cross_dock varchar(50),
		 vendor_commission NUMERIC(10, 2),
		 vendor_remark varchar(50),
		 vendor_base_rate NUMERIC(10, 2),
		 vendor_km_cost NUMERIC(10, 2),
		 vendor_toll NUMERIC(10, 2),
		 vendor_total_hire NUMERIC(10, 2),
		 vendor_advance NUMERIC(10, 2),
		 vendor_debit_amount NUMERIC(10, 2),
		 vendor_balance_amount NUMERIC(10, 2),
		 vendor_break_down VARCHAR(255),

		 vendor_paid_by varchar(100),
		 vendor_load_unload_amount NUMERIC(10, 2),
		 vendor_halting_days NUMERIC(10, 2),
		 vendor_halting_paid NUMERIC(10, 2),
		 vendor_extra_delivery NUMERIC(10, 2),

		 vendor_monul NUMERIC(10, 2),
		 vendor_total_amount NUMERIC(10, 2),
		 
		 load_status varchar(30),

		 pod_received INT,
		 CONSTRAINT tp_pod_received_0_or_1 CHECK (pod_received IN (0, 1)),
		 is_lr_generated BOOLEAN DEFAULT FALSE,
    	 CONSTRAINT tp_lr_gen_0_or_1 CHECK (is_lr_generated IN (0, 1)),

		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (trip_sheet_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = tripSheet.Exec()
	if err != nil {
		lg.Println("ERROR: trip_sheet_header table: ", err.Error())
	}
}

func createTripSheetHeaderLoadingPoints(db *sql.DB) {
	tripSheet, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		trip_sheet_header_load_unload_points (load_unload_point_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		trip_sheet_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_trip_sheet_load_ref FOREIGN KEY (trip_sheet_id) REFERENCES trip_sheet_header(trip_sheet_id),

		loading_point_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_tp_load_ref FOREIGN KEY (loading_point_id) REFERENCES loading_point (loading_point_id),
		
		type varchar(100) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (load_unload_point_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = tripSheet.Exec()
	if err != nil {
		lg.Println("ERROR: trip_sheet_header table: ", err.Error())
	}
}

func createManagePod(db *sql.DB) {
	managePod, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
		manage_pod (pod_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		trip_sheet_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_manage_pod_ref FOREIGN KEY (trip_sheet_id) REFERENCES trip_sheet_header(trip_sheet_id),
		trip_sheet_num varchar(100) NOT NULL,
		CONSTRAINT fky_manage_pod_tpnum_ref FOREIGN KEY (trip_sheet_num) REFERENCES trip_sheet_header(trip_sheet_num),
		lr_number varchar(50) UNIQUE NOT NULL,
		customer_name varchar(50),
		customer_id INT UNSIGNED,
		CONSTRAINT fky_pod_cus_ref FOREIGN KEY (customer_id) REFERENCES customers (customer_id),
		send_by varchar(100),
		pod_status varchar(100),
		pod_remark varchar(1000),
		late_submission_debit DECIMAL(10,2),
		pod_submited_date varchar(50),
		halting_amount DECIMAL(10,2),
		paid_by varchar(50),
		unloading_date varchar(50),
        unloading_charges DECIMAL(10,2),
		pod_doc varchar(1000),
		trip_type varchar(50),
		kilometers_covered DECIMAL(10,2),
		halting_days varchar(50),
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_pod_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (pod_id)
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = managePod.Exec()
	if err != nil {
		lg.Println("ERROR: pod_id table: ", err.Error())
	}
}

func createLRHeader(db *sql.DB) {

	lrReceipt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS 
			lr_receipt_header (lr_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			trip_sheet_id INT UNSIGNED NOT NULL,
			CONSTRAINT fky_lr_tsid_ref FOREIGN KEY (trip_sheet_id) REFERENCES trip_sheet_header(trip_sheet_id),
			trip_sheet_num varchar(100) NOT NULL,
			lr_number varchar(50) UNIQUE NOT NULL,
			trip_date varchar(50),
			vehicle_number varchar(50) NOT NULL,
			vehicle_size varchar(100),
			driver_name varchar(50) NOT NULL,
			driver_mobile_number varchar(20) NOT NULL,
			invoice_number varchar(100) NOT NULL,
			invoice_value DECIMAL(10,2) NOT NULL,
			consignor_name varchar(50) NOT NULL,
			consignor_address varchar(1000) NOT NULL,
			consignor_gst varchar(50),
			consignee_name varchar(50) NOT NULL,
			consignee_address varchar(1000) NOT NULL,
			consignee_gst varchar(50),
			goods_type varchar(50),
			goods_weight varchar(50),
			quantity_in_pieces varchar(50),
			remark varchar(1000),
			org_id INT UNSIGNED NOT NULL,
			CONSTRAINT fky_lr_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (lr_id)
			) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = lrReceipt.Exec()
	if err != nil {
		lg.Println("ERROR: lrReceipt table: ", err.Error())
	}
}

func createCustomers(db *sql.DB) {
	customer, err := db.Prepare(`CREATE TABLE IF NOT EXISTS customers (
    customer_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	customer_code VARCHAR(50) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
	branch_name VARCHAR(100),
	branch_id INT UNSIGNED,
	nick_name VARCHAR(50),
	payment_terms VARCHAR(50),
	gst_number VARCHAR(50),
	remark VARCHAR(255),
	kkt_responsible_emp VARCHAR(50),
	kkt_responsible_emp_id INT UNSIGNED,
	CONSTRAINT fky_kkt_responsible_emp_org FOREIGN KEY (kkt_responsible_emp_id) REFERENCES employee (emp_id),
	fassi VARCHAR(50),
	pan_number VARCHAR(50),
	address VARCHAR(1000),
	city VARCHAR(100),
	state VARCHAR(100),
	pincode INT UNSIGNED,
	country varchar(100),
	status VARCHAR(20),
	is_active BOOLEAN DEFAULT TRUE,
	org_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_customers_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
	company_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = customer.Exec()
	if err != nil {
		lg.Println("ERROR: customers table: ", err.Error())
	}
}

func createCustomersContactInfo(db *sql.DB) {
	contactInfo, err := db.Prepare(`CREATE TABLE IF NOT EXISTS contact_info (
    contact_info_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	customer_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_customer_ct_info FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
	contact_person_name VARCHAR(50) NOT NULL,
    post VARCHAR(50),
	email_id VARCHAR(100),
	contact_nummber1 VARCHAR(50),
	contact_nummber2 VARCHAR(50),
	contact_nummber3 VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = contactInfo.Exec()
	if err != nil {
		lg.Println("ERROR: contact_info table: ", err.Error())
	}
}
func createCustomersAggrement(db *sql.DB) {
	customer, err := db.Prepare(`CREATE TABLE IF NOT EXISTS customer_aggrement (
    aggrement_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	customer_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_customer_agg FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
	aggrement_number VARCHAR(10),
	aggrement_period VARCHAR(50),
	aggrement_name VARCHAR(100),
	aggrement_type VARCHAR(100),
	aggrement_doc VARCHAR(1000),
    remark VARCHAR(1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = customer.Exec()
	if err != nil {
		lg.Println("ERROR: customer_aggrement table: ", err.Error())
	}
}

func createVendors(db *sql.DB) {
	vendor, err := db.Prepare(`CREATE TABLE IF NOT EXISTS vendors (
		vendor_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		vendor_name VARCHAR(100) NOT NULL,
		vendor_code VARCHAR(50) NOT NULL,
		owner_name VARCHAR(50),
		gst_number VARCHAR(50),
		preferred_operating_routes VARCHAR(100),
		address_line1 TEXT,
		city VARCHAR(100),
		state VARCHAR(100),
		pan_number VARCHAR(50),
		tds_declaration VARCHAR(100),
		remark TEXT,
		bank_account_holder_name VARCHAR(100),
		bank_account_number VARCHAR(100),
		bank_name VARCHAR(100),
		bank_ifsc_code VARCHAR(20),
		pancard_img VARCHAR(1000),
		bank_passbook_or_cheque_img VARCHAR(1000),
		status VARCHAR(20) NOT NULL,
		is_active BOOLEAN DEFAULT TRUE NOT NULL,
		org_id INT UNSIGNED NOT NULL,
		CONSTRAINT fky_vendors_org FOREIGN KEY (org_id) REFERENCES organisation(org_id),
		login_type VARCHAR(10),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (vendor_id),
		CONSTRAINT vendorx_vendors_code_exists UNIQUE (vendor_code)
	)ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = vendor.Exec()
	if err != nil {
		lg.Println("ERROR: vendors table: ", err.Error())
	}
}

func createDriverContactInfo(db *sql.DB) {
	contactInfo, err := db.Prepare(`CREATE TABLE IF NOT EXISTS vendor_contact_info (
    contact_info_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	vendor_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_vendor_ct_info FOREIGN KEY (vendor_id) REFERENCES vendors(vendor_id),
	contact_person_name VARCHAR(50) NOT NULL,
    post VARCHAR(50),
	email_id VARCHAR(100),
	contact_nummber1 VARCHAR(50),
	contact_nummber2 VARCHAR(50),
	contact_nummber3 VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = contactInfo.Exec()
	if err != nil {
		lg.Println("ERROR: vendor_contact_info table: ", err.Error())
	}
}

func createDeclarationYearInfo(db *sql.DB) {
	contactInfo, err := db.Prepare(`CREATE TABLE IF NOT EXISTS declaration_year_info (
    declaration_year_info_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	vendor_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_vendor_dec_info FOREIGN KEY (vendor_id) REFERENCES vendors(vendor_id),
	declaration_year VARCHAR(50),
	declaration_doc VARCHAR(1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = contactInfo.Exec()
	if err != nil {
		lg.Println("ERROR: declaration_year_info table: ", err.Error())
	}
}

//

func createAppConfig(db *sql.DB) {
	contactInfo, err := db.Prepare(`CREATE TABLE IF NOT EXISTS app_config (
    app_config_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	config_code VARCHAR(50) NOT NULL,
	config_name VARCHAR(50) NOT NULL,
	value VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = contactInfo.Exec()
	if err != nil {
		lg.Println("ERROR: app_config table: ", err.Error())
	}
}

func createEmployeeAttendance(db *sql.DB) {

	attendanceTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS attendance (
    attendance_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	employee_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_employee_attendance FOREIGN KEY (employee_id) REFERENCES employee (emp_id),
    check_in_date DATETIME NOT NULL,
    check_in_date_str VARCHAR(50) NOT NULL,
    check_out_date DATETIME,
    check_out_date_str VARCHAR(50),
    in_time VARCHAR(20) NOT NULL,
    out_time VARCHAR(20),
    late_time INT DEFAULT 0,
    over_time INT DEFAULT 0,
    duration VARCHAR(20),
    hours INT DEFAULT 0,
    minutes INT DEFAULT 0,
    check_in_latitude DOUBLE,
    check_in_longitude DOUBLE,
    check_in_area VARCHAR(100),
    check_in_city VARCHAR(100),
    check_out_latitude DOUBLE,
    check_out_longitude DOUBLE,
    check_out_area VARCHAR(100),
    check_out_city VARCHAR(100),
    is_completed BOOLEAN DEFAULT FALSE,
	manual_checked_in BOOLEAN DEFAULT FALSE,
    manual_checked_out BOOLEAN DEFAULT FALSE,
    check_out_by_id BIGINT,
    check_in_by_id BIGINT,
    status VARCHAR(50),
	CONSTRAINT attendance_emp_date_unique UNIQUE (employee_id, check_in_date_str),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = attendanceTable.Exec()
	if err != nil {
		lg.Println("ERROR: attendance table: ", err.Error())
	}

}

func createEmployeeAttendanceEntry(db *sql.DB) {

	attendanceTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS attendance_entry_log (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	attendance_id BIGINT UNSIGNED NOT NULL,
	CONSTRAINT fky_attendance_ref FOREIGN KEY (attendance_id) REFERENCES attendance (attendance_id),
	employee_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_emp_att_id FOREIGN KEY (employee_id) REFERENCES employee (emp_id),
    entry_in_date DATETIME NOT NULL,
    entry_in_date_str VARCHAR(50) NOT NULL,
	entry_time VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = attendanceTable.Exec()
	if err != nil {
		lg.Println("ERROR: attendance_entry_log table: ", err.Error())
	}

}

func createCustomerInvoice(db *sql.DB) {

	customerInvoice, err := db.Prepare(`CREATE TABLE IF NOT EXISTS customer_invoice (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		invoice_number VARCHAR(100),
		invoice_ref VARCHAR(50) NOT NULL,
		work_type VARCHAR(50) NOT NULL,
		work_start_date VARCHAR(50) NOT NULL,
		work_end_date VARCHAR(50) NOT NULL,
		document_date VARCHAR(50) NOT NULL,
		invoice_status VARCHAR(50) NOT NULL,
		invoice_amount DOUBLE NOT NULL,
		invoice_date VARCHAR(50),
		payment_date VARCHAR(50),
		customer_id INT UNSIGNED,
		customer_name VARCHAR(100),
		customer_code VARCHAR(100),
		transaction_id VARCHAR(255),
		trip_ref LONGTEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = customerInvoice.Exec()
	if err != nil {
		lg.Println("ERROR: customerInvoice table: ", err.Error())
	}
}

func createVehicles(db *sql.DB) {
	vehicle, err := db.Prepare(`CREATE TABLE IF NOT EXISTS vehicles (
    vehicle_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	vendor_id INT UNSIGNED NOT NULL,
	CONSTRAINT fky_vendors_ref_info FOREIGN KEY (vendor_id) REFERENCES vendors(vendor_id),
	vehicle_number VARCHAR(50) NOT NULL,
    vehicle_type VARCHAR(50) NOT NULL,
	vehicle_make VARCHAR(50),
	vehicle_model VARCHAR(50) NOT NULL,
	permit_type VARCHAR(50),
	vehicle_size VARCHAR(50),
	closed_open VARCHAR(50),
	vehicle_capacity_tons VARCHAR(50),
	rc_expiry_doc VARCHAR(500),
	insurance_doc VARCHAR(500),
	pucc_expiry_doc VARCHAR(500),
	np_expire_doc VARCHAR(500),
	fitness_expiry_doc VARCHAR(500),
	tax_expiry_doc VARCHAR(500),
	mp_expire_doc VARCHAR(500),
	CONSTRAINT vendor_vehicle_unique UNIQUE (vendor_id, vehicle_number),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = vehicle.Exec()
	if err != nil {
		lg.Println("ERROR: vendor_vehiclecontact_info table: ", err.Error())
	}
}

func createVehicleSizeTypes(db *sql.DB) {
	contactInfo, err := db.Prepare(`CREATE TABLE IF NOT EXISTS vehicle_size_type (
    vehicle_size_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    vehicle_size VARCHAR(100),
	vehicle_type VARCHAR(50),
	status VARCHAR(50),
	is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )ENGINE=InnoDB;`)
	if err != nil {
		lg.Println(err.Error())
	}
	_, err = contactInfo.Exec()
	if err != nil {
		lg.Println("ERROR: vehicle_size_type table: ", err.Error())
	}
}

func insertSuperAdmin(db *sql.DB) {
	organisation := schema.Organisation{}
	row := db.QueryRow("SELECT * FROM organisation LIMIT 1;")
	err := row.Scan(&organisation.OrgId)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	lg.Print("insertSuperAdmin..........")
	organisation.Name = "KK Transport"
	organisation.DisplayName = "KK Transport"
	organisation.DomainName = "kkt.in"
	organisation.EmailId = "prabhasjraj@gmail.com"
	organisation.ContactName = "KKT Admin"
	organisation.ContactNo = "9876543210"
	organisation.City = "Chennai"
	organisation.Version = 1

	query := fmt.Sprintf(`INSERT INTO 
		organisation (
		name, 
		display_name, 
		domain_name, 
		email_id, 
		contact_name, 
		contact_no, 
		city,
		version) 
		values('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')`,
		organisation.Name,
		organisation.DisplayName,
		organisation.DomainName,
		organisation.EmailId,
		organisation.ContactName,
		organisation.ContactNo,
		organisation.City,
		organisation.Version)

	result, err := db.Exec(query)
	if err != nil {
		lg.Println("Error db.Exec(query) org: ", err)
		return
	}
	orgID, err := result.LastInsertId()
	if err != nil {
		lg.Println("Error getting last insert ID: ", err)
	}
	fmt.Println("First Organisation ID: ", orgID)

	newRole := schema.Role{
		RoleCode:    "SUPER_ADMIN",
		RoleName:    "Super Administrator",
		Description: "Has full access to all resources.",
		OrgId:       orgID,
		IsActive:    1,
		Version:     1,
	}
	roleQuery := fmt.Sprintf(`INSERT INTO thub_role (
		role_code, 
		role_name, 
		description, 
		org_id, 
		is_active, 
		version) VALUES ('%v', '%v', '%v', %v, %v, %v)`,
		newRole.RoleCode,
		newRole.RoleName,
		newRole.Description,
		newRole.OrgId,
		newRole.IsActive,
		newRole.Version)

	roleResult, err := db.Exec(roleQuery)
	if err != nil {
		lg.Println("Error db.Exec(userQuery) user_login: ", err)
		return
	}
	roleId, err := roleResult.LastInsertId()
	if err != nil {
		lg.Println("Error getting last insert ID: ", err)
		return
	}

	user := &schema.UserLogin{
		FirstName:    organisation.Name,
		LastName:     organisation.DisplayName,
		EmailId:      organisation.EmailId,
		MobileNo:     organisation.ContactNo,
		Password:     utils.SHAEncoding("welcome123"),
		RoleID:       roleId,
		OrgId:        orgID,
		IsSuperAdmin: 1,
		LoginType:    "Web",
	}
	userQuery := fmt.Sprintf(`INSERT INTO 
		user_login ( first_name, 
		last_name, 
		email_id, 
		mobile_no, 
		password, 
		role_id, 
		org_id,
		is_super_admin,
		login_type,
		version) 
		values('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')`,
		user.FirstName,
		user.LastName,
		user.EmailId,
		user.MobileNo,
		user.Password,
		user.RoleID,
		user.OrgId,
		user.IsSuperAdmin,
		user.LoginType,
		organisation.Version)

	userResult, err := db.Exec(userQuery)
	if err != nil {
		lg.Println("Error db.Exec(userQuery) user_login: ", err)
		return
	}
	userID, err := userResult.LastInsertId()
	if err != nil {
		lg.Println("Error getting last insert ID: ", err)
	}
	fmt.Println("First user: ", userID, user.FirstName)
}

func insertTruckTypes(db *sql.DB) {
	trucktId := schema.TruckTypes{}
	row := db.QueryRow("SELECT * FROM truck_types LIMIT 1;")
	err := row.Scan(&trucktId.TruckTypeId)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	lg.Print("insertTruckTypes..........")
	truckTypes := []struct {
		Category  string `json:"category"`
		TruckType string `json:"truck_type"`
	}{
		// {"Open", "17 Feet Open"},
		// {"Open", "20 Feet Open"},
		// {"Open", "22 Feet Open"},
		// {"Open", "24 Feet Open"},
		// {"Open", "10 Whl Open"},
		// {"Open", "12 Whl Open"},
		// {"Open", "14 Whl Open"},
		// {"Closed", "32 Feet Single Axle"},
		// {"Closed", "32 Feet Single Axle High Cube"},
		// {"Closed", "32 Feet Multi Axle"},
		// {"Closed", "32 Feet Multi Axle High Cube"},
		// {"Closed", "32 Feet Triple Axle"},
		// {"Closed", "20 Feet Closed"},
		// {"Closed", "22 Feet Closed"},
		// {"Closed", "24 Feet Closed"},

		{"Closed", "Dost"},
		{"Closed", "10 Feet"},
		{"Closed", "14 Feet Multi Axle"},
		{"Closed", "17 Feet Multi Axle High Cube"},
		{"Closed", "20 Feet Triple Axle"},
		{"Closed", "22 Feet Closed"},
		{"Closed", "24 Feet Closed"},
		{"Closed", "32 Feet sxl Single Axle"},
		{"Closed", "32 Feet mxl (Multi Axle)"},
	}

	for _, truck := range truckTypes {
		query := fmt.Sprintf(`INSERT INTO truck_types (
			category, 
			truck_type, 
			is_active, 
			org_id) 
			values('%v', '%v', %v, %v)`,
			truck.Category,
			truck.TruckType,
			true,
			1)

		result, err := db.Exec(query)
		if err != nil {
			lg.Println("Error inserting truck type:", err)
			continue
		}

		_, errM := result.LastInsertId()
		if errM != nil {
			lg.Println("Error getting last insert ID for truck type:", errM)
			continue
		}

		//fmt.Printf("Inserted Truck Type ID: %d, Type: %s\n", truckTypeID, truck.TruckType)
	}
}

func insertAppConfig(db *sql.DB) {
	appConfigID := schema.AppConfig{}
	row := db.QueryRow("SELECT * FROM app_config LIMIT 1;")
	err := row.Scan(&appConfigID.AppConfigID)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	lg.Print("insertTruckTypes..........")
	appConfig := []struct {
		ConfigCode  string `json:"config_code"`
		ConfigName  string `json:"config_name"`
		ConfigValue string `json:"value"`
	}{
		{constant.SESSION_TIMEOUT, constant.SESSION_TIMEOUT_MS, constant.SESSION_TIMEOUT_VALUE},
		{constant.EMPLOYEE_PERFORMANCE, constant.EMPLOYEE_PERFORMANCE_1, constant.EMPLOYEE_PERFORMANCE_1_NAME},
		{constant.EMPLOYEE_PERFORMANCE, constant.EMPLOYEE_PERFORMANCE_2, constant.EMPLOYEE_PERFORMANCE_2_NAME},
		{constant.EMPLOYEE_PERFORMANCE, constant.EMPLOYEE_PERFORMANCE_3, constant.EMPLOYEE_PERFORMANCE_3_NAME},
		{constant.EMPLOYEE_PERFORMANCE, constant.EMPLOYEE_PERFORMANCE_4, constant.EMPLOYEE_PERFORMANCE_4_NAME},
		{constant.EMPLOYEE_PERFORMANCE, constant.EMPLOYEE_PERFORMANCE_5, constant.EMPLOYEE_PERFORMANCE_5_NAME},
	}

	for _, config := range appConfig {
		query := fmt.Sprintf(`INSERT INTO app_config (
			config_code, 
			config_name, 
			value) 
			values('%v', '%v', '%v')`,
			config.ConfigCode,
			config.ConfigName,
			config.ConfigValue)

		result, err := db.Exec(query)
		if err != nil {
			lg.Println("Error inserting app_config:", err)
			continue
		}

		_, errM := result.LastInsertId()
		if errM != nil {
			lg.Println("Error getting last insert ID for app_config:", errM)
			continue
		}
	}
}

func insertMaterialTypes(db *sql.DB) {
	trucktId := schema.TruckTypes{}
	row := db.QueryRow("SELECT * FROM material_types LIMIT 1;")
	err := row.Scan(&trucktId.TruckTypeId)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	materialTypes := []struct {
		MaterialName string `json:"material_name"`
	}{
		{"Agriculture"},
		{"Alcoholic Beverage"},
		{"Auto Parts"},
		{"Bags"},
		{"Beverages"},
		{"Bundle / Rolls"},
		{"Carton Box"},
		{"Carton box / FMCG"},
		{"Drum/barrel"},
		{"Electronic Goods"},
		{"Empty Bottle"},
		{"Empty Cylinder"},
		{"Engineering materials"},
		{"FMCG"},
		{"Footwears"},
		{"Furniture"},
		{"Garment bundle"},
		{"Glass"},
		{"Hardware (Door / Windows)"},
		{"Hardware (Pipes / Tiles)"},
		{"Industrial Goods"},
		{"Kitchenware"},
		{"Medicine"},
		{"Milk Products"},
		{"Packaging Materials"},
		{"Paint"},
		{"Pallet / Panel"},
		{"Paper Bundle"},
		{"Parcel"},
		{"Plywood"},
		{"Sack (Jute/Plastic)"},
		{"Scrap"},
		{"Solid Sheets"},
		{"Textiles"},
		{"Tyre"},
		{"Waste Cloths"},
		{"Waste Material"},
		{"Others"},
	}

	for _, material := range materialTypes {
		query := fmt.Sprintf(`INSERT INTO material_types (
			material_name,
			is_active,
			org_id)
			VALUES ('%s', %v, %v)`,
			material.MaterialName, true, 1)

		result, err := db.Exec(query)
		if err != nil {
			lg.Println("Error inserting material type:", err)
			continue
		}

		_, errMa := result.LastInsertId()
		if errMa != nil {
			lg.Println("Error getting last insert ID for material type:", errMa)
			continue
		}
		//fmt.Printf("Inserted Material Type ID: %d, Name: %s\n", materialTypeID, material.MaterialName)
	}
}

func insertRolePermisstionLabel(db *sql.DB) {

	permisstionLabel := schema.PermissionLabel{}
	row := db.QueryRow("SELECT * FROM permisstion_label LIMIT 1;")
	err := row.Scan(&permisstionLabel.PermisstionLabel)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	lg.Print("insertRolePermisstionLabel..........")

	permisstionLabels := []struct {
		PermisstionLabel string `json:"permisstion_label"`
		Description      string `json:"description"`
	}{
		{"EDIT", "edit access"},
		{"VIEW", "view access"},
		{"No", "No access"},
	}

	for _, permisstion := range permisstionLabels {
		query := fmt.Sprintf(`INSERT INTO permisstion_label (
			permisstion_label, description ) values('%v', '%v')`,
			permisstion.PermisstionLabel,
			permisstion.Description)

		lg.Println("inserting permisstionLabels:", query)
		result, err := db.Exec(query)
		if err != nil {
			lg.Println("Error inserting permisstionLabels:", err)
			continue
		}

		_, errM := result.LastInsertId()
		if errM != nil {
			lg.Println("Error getting last insert ID for permisstionLabels:", errM)
			continue
		}

		//fmt.Printf("Inserted Truck Type ID: %d, Type: %s\n", truckTypeID, truck.TruckType)
	}
}

func insertwebsiteScreens(db *sql.DB) {

	webScreen := schema.WebsiteScreens{}
	row := db.QueryRow("SELECT * FROM website_screens LIMIT 1;")
	err := row.Scan(&webScreen.ScreenName)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	lg.Print("insertwebsiteScreens..........")

	webScreens := []struct {
		ScreenName  string `json:"screen_name"`
		Description string `json:"description"`
	}{
		{constant.MENU_1_OVERVIEW, constant.MENU_1_OVERVIEW_DESC},
		{constant.MENU_2_TRIP_MANAGEMENT, constant.MENU_2_TRIP_MANAGEMENT_DESC},
		{constant.MENU_3_VENDORS, constant.MENU_3_VENDORS_DESC},
		{constant.MENU_4_CUSTOMERS, constant.MENU_4_CUSTOMERS_DESC},
		{constant.MENU_5_OPERATIONS, constant.MENU_5_OPERATIONS_DESC},
		{constant.MENU_6_SETTINGS, constant.MENU_6_SETTINGS_DESC},
		{constant.MENU_7_REPORTS, constant.MENU_7_REPORTS_DESC},
		{constant.MENU_8_EMPLOYEE, constant.MENU_8_EMPLOYEE_DESC},
	}

	for _, screen := range webScreens {
		query := fmt.Sprintf(`INSERT INTO website_screens (
			screen_name, description ) values('%v', '%v')`,
			screen.ScreenName,
			screen.Description)

		lg.Println("inserting webScreens:", query)
		result, err := db.Exec(query)
		if err != nil {
			lg.Println("Error inserting webScreens:", err)
			continue
		}

		_, errM := result.LastInsertId()
		if errM != nil {
			lg.Println("Error getting last insert ID for webScreens:", errM)
			continue
		}
		//fmt.Printf("Inserted Truck Type ID: %d, Type: %s\n", truckTypeID, truck.TruckType)
	}
}

// func AboutPrivacyTableCreation(db *sql.DB) {
// 	aboutPrivacyTable, err := db.Prepare(`CREATE TABLE IF NOT EXISTS
// 		ac_about_privacy_policy (about_privacy_policy_id int unsigned NOT NULL AUTO_INCREMENT,
// 		version_code int unsigned NOT NULL,
// 		version_name varchar(255) NOT NULL,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 		CONSTRAINT versionx_about_privacy_policy UNIQUE (version_code),
// 		PRIMARY KEY (about_privacy_policy_id));`)
// 	if err != nil {
// 		lg.Println(err.Error())
// 	}
// 	_, err = aboutPrivacyTable.Exec()
// 	if err != nil {
// 		lg.Println(err.Error())
// 	}
// 	aboutPrivacyInfoTable, errN := db.Prepare(`CREATE TABLE IF NOT EXISTS
// 		ac_about_privacy_policy_info (id int unsigned NOT NULL AUTO_INCREMENT,
// 	    about_privacy_policy_id int unsigned NOT NULL,
// 		CONSTRAINT fk_ac_about_privacy_policy FOREIGN KEY (about_privacy_policy_id) REFERENCES ac_about_privacy_policy(about_privacy_policy_id),
// 		version_code int unsigned NOT NULL,
// 		version_name varchar(255) NOT NULL,
// 		sequence_no int unsigned NOT NULL,
// 		content_type varchar(1000) NOT NULL,
// 		message_info BLOB,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 		CONSTRAINT versionx_about_privacy_policy_info UNIQUE (version_code,about_privacy_policy_id,sequence_no),
// 		PRIMARY KEY (id));`)
// 	if errN != nil {
// 		lg.Println(errN.Error())
// 	}
// 	_, err = aboutPrivacyInfoTable.Exec()
// 	if err != nil {
// 		lg.Println(err.Error())
// 	}
// }

// func BehaviourChangeTechniquesTableCreation(db *sql.DB) {
// 	ac_bct, err := db.Prepare(`CREATE TABLE IF NOT EXISTS ac_behaviour_change_techniques (behaviour_change_id int unsigned NOT NULL AUTO_INCREMENT,
// 		bct_taxonomy_id varchar(100) NOT NULL,
// 		bct_taxonomy varchar(255) NOT NULL,
// 		bct_id varchar(100) NOT NULL,
// 		bct_description varchar(1000) NOT NULL,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 		CONSTRAINT bct_idx_behaviour_change UNIQUE (bct_id),
// 		PRIMARY KEY (behaviour_change_id));`)
// 	if err != nil {
// 		lg.Println(err.Error())
// 	}
// 	_, err = ac_bct.Exec()
// 	if err != nil {
// 		lg.Println("Error -", err.Error())
// 	}

// }

// func BehaviourChangeInterventionsTableCreation(db *sql.DB) {
// 	ac_bct, err := db.Prepare(`CREATE TABLE IF NOT EXISTS ac_behaviour_change_notifications (behaviour_notification_id int unsigned NOT NULL AUTO_INCREMENT,
// 		bcn_category varchar(100) NOT NULL,
// 		bcn_category_desc varchar(255),
// 		bcn_group varchar(100),
// 		bcn_group_description varchar(255),
// 		bcn_trigger_event varchar(255),
// 		app_route varchar(255),
// 		bcn_id varchar(100) NOT NULL,
// 		bcn_message varchar(2555),
// 		bct_1 varchar(100),
// 		bct_2 varchar(100),
// 		bct_3 varchar(100),
// 		bct_4 varchar(100),
// 		alcochange_theme BOOLEAN DEFAULT FALSE,
// 		frames_feedback BOOLEAN DEFAULT FALSE,
// 		frames_responsibility BOOLEAN DEFAULT FALSE,
// 		frames_advice BOOLEAN DEFAULT FALSE,
// 		frames_menu BOOLEAN DEFAULT FALSE,
// 		frames_empathy BOOLEAN DEFAULT FALSE,
// 		frames_support_and_selfefficacy BOOLEAN DEFAULT FALSE,
// 		develop_discrepancy BOOLEAN DEFAULT FALSE,
// 		assessment BOOLEAN DEFAULT FALSE,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 		CONSTRAINT bcn_idx_behaviour_change UNIQUE (bcn_id),
// 		PRIMARY KEY (behaviour_notification_id));`)
// 	if err != nil {
// 		lg.Println(err.Error())
// 	}
// 	_, err = ac_bct.Exec()
// 	if err != nil {
// 		lg.Println("Error -", err.Error())
// 	}

// }

// func PatientEngagementReminderTableCreation(db *sql.DB) {
// 	acIt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS ac_intervention_type (id int unsigned NOT NULL AUTO_INCREMENT,
// 		intervention_type_id int unsigned NOT NULL UNIQUE,
// 		intervention_type varchar(255) NOT NULL,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 		PRIMARY KEY (id, intervention_type_id));`)
// 	if err != nil {
// 		lg.Println(err.Error())
// 	}
// 	_, err = acIt.Exec()
// 	if err != nil {
// 		lg.Println("Error -", err.Error())
// 	}

// 	patientEngagementReminderTable, errN := db.Prepare(`CREATE TABLE IF NOT EXISTS
// 		ac_patient_engagement_reminder (id int unsigned NOT NULL AUTO_INCREMENT,
// 		user_id bigint unsigned NOT NULL,
// 		CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user(user_id),
// 		user_uuid varchar(100) NOT NULL,
// 	    intervention_type_id int unsigned NOT NULL,
// 		CONSTRAINT fk_ac_intervention_type FOREIGN KEY (intervention_type_id) REFERENCES ac_intervention_type(intervention_type_id),
// 		notification_id int unsigned NOT NULL,
// 		patient_engagement_time varchar(100) NOT NULL,
// 		message_shown varchar(1000),
// 		user_action varchar(100) NOT NULL,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 		CONSTRAINT reminderx_patient_engagement UNIQUE (notification_id,patient_engagement_time),
// 		PRIMARY KEY (id));`)
// 	if errN != nil {
// 		lg.Println(errN.Error())
// 	}
// 	_, err = patientEngagementReminderTable.Exec()
// 	if err != nil {
// 		lg.Println(err.Error())
// 	}

// }

// func BenefitTherapy(db *sql.DB) {

// 	ac_bt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS ac_benefit_therapy (id int unsigned NOT NULL AUTO_INCREMENT,
// 		message varchar(255) NOT NULL,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 		PRIMARY KEY (id));`)
// 	if err != nil {
// 		lg.Println(err.Error())
// 	}
// 	_, err = ac_bt.Exec()
// 	if err != nil {
// 		lg.Println("Error -", err.Error())
// 	}
// }
