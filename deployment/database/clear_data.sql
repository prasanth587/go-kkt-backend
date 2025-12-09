SET FOREIGN_KEY_CHECKS = 0;

-- Transactional Tables
TRUNCATE TABLE trip_sheet_header;
TRUNCATE TABLE trip_sheet_header_load_unload_points;
TRUNCATE TABLE trip_sheet_num;
TRUNCATE TABLE attendance;
TRUNCATE TABLE attendance_entry_log;
TRUNCATE TABLE lr_receipt_header;
TRUNCATE TABLE manage_pod;
TRUNCATE TABLE customer_invoice;

-- Master Data Tables
TRUNCATE TABLE vendors;
TRUNCATE TABLE vendor;
TRUNCATE TABLE vendor_contact_info;

TRUNCATE TABLE vehicles;
TRUNCATE TABLE vehicle;
TRUNCATE TABLE declaration_year_info;
TRUNCATE TABLE customers;
TRUNCATE TABLE customer;
TRUNCATE TABLE contact_info;
TRUNCATE TABLE customer_aggrement;
TRUNCATE TABLE loading_point;
TRUNCATE TABLE branch;
TRUNCATE TABLE driver;
TRUNCATE TABLE employee;

-- CAUTION: Users and Roles
-- Uncomment the following lines if you want to clear users and roles too.
-- NOTE: This will delete your ADMIN account!
-- TRUNCATE TABLE user_login;
-- TRUNCATE TABLE thub_role;
-- TRUNCATE TABLE thub_role_screens_map;

SET FOREIGN_KEY_CHECKS = 1;



