-- WARNING: This will DELETE ALL customer data!
-- Make sure you have a backup before running this!

SET FOREIGN_KEY_CHECKS = 0;

-- Delete related data first (child tables)
DELETE FROM trip_sheet_header WHERE customer_id IS NOT NULL;
DELETE FROM customer_invoice WHERE customer_id IS NOT NULL;
DELETE FROM contact_info;
DELETE FROM customer_aggrement;

-- Now delete customers
DELETE FROM customers;

-- Reset auto-increment to start from 1
ALTER TABLE customers AUTO_INCREMENT = 1;

SET FOREIGN_KEY_CHECKS = 1;

-- Verify it's empty
SELECT COUNT(*) as remaining_customers FROM customers;
SELECT 'Customers deleted successfully! Next customer will start from CUS-000001' as message;


