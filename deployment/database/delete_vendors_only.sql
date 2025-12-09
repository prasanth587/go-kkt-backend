-- WARNING: This will DELETE ALL vendor data!
-- Make sure you have a backup before running this!

SET FOREIGN_KEY_CHECKS = 0;

-- Delete related data first (child tables)
DELETE FROM trip_sheet_header WHERE vendor_id IS NOT NULL;
DELETE FROM vehicles WHERE vendor_id IS NOT NULL;
DELETE FROM vendor_contact_info;
DELETE FROM declaration_year_info;

-- Now delete vendors
DELETE FROM vendors;

-- Reset auto-increment to start from 1
ALTER TABLE vendors AUTO_INCREMENT = 1;

SET FOREIGN_KEY_CHECKS = 1;

-- Verify it's empty
SELECT COUNT(*) as remaining_vendors FROM vendors;
SELECT 'Vendors deleted successfully! Next vendor will start from VEN-000001' as message;


