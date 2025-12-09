-- WARNING: This will DELETE ALL vendor AND customer data!
-- Make sure you have a backup before running this!

SET FOREIGN_KEY_CHECKS = 0;

-- Delete trip sheet data (references both vendors and customers)
DELETE FROM trip_sheet_header_load_unload_points;
DELETE FROM trip_sheet_header;

-- Delete vendor related data
DELETE FROM vehicles;
DELETE FROM vendor_contact_info;
DELETE FROM declaration_year_info;
DELETE FROM vendors;

-- Delete customer related data
DELETE FROM customer_invoice;
DELETE FROM contact_info;
DELETE FROM customer_aggrement;
DELETE FROM customers;

-- Reset auto-increment counters
ALTER TABLE vendors AUTO_INCREMENT = 1;
ALTER TABLE customers AUTO_INCREMENT = 1;

SET FOREIGN_KEY_CHECKS = 1;

-- Verify
SELECT 
    (SELECT COUNT(*) FROM vendors) as remaining_vendors,
    (SELECT COUNT(*) FROM customers) as remaining_customers;
SELECT 'Vendors and Customers deleted successfully! They will start from 000001' as message;


