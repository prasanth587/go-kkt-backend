-- Delete trip sheet data for IDs 6, 7, and 8
-- Also delete vendor 3 and all related data
-- This script deletes from related tables first, then from the main tables

SET FOREIGN_KEY_CHECKS = 0;

-- Delete from related tables (order matters due to foreign key constraints)

-- 1. Delete load/unload points for these trip sheets
DELETE FROM trip_sheet_header_load_unload_points WHERE trip_sheet_id = 6;
DELETE FROM trip_sheet_header_load_unload_points WHERE trip_sheet_id = 7;
DELETE FROM trip_sheet_header_load_unload_points WHERE trip_sheet_id = 8;

-- 2. Delete POD records for these trip sheets
DELETE FROM manage_pod WHERE trip_sheet_id = 6;
DELETE FROM manage_pod WHERE trip_sheet_id = 7;
DELETE FROM manage_pod WHERE trip_sheet_id = 8;

-- 3. Delete LR receipt records for these trip sheets
DELETE FROM lr_receipt_header WHERE trip_sheet_id = 6;
DELETE FROM lr_receipt_header WHERE trip_sheet_id = 7;
DELETE FROM lr_receipt_header WHERE trip_sheet_id = 8;

-- 4. Finally, delete from the main trip_sheet_header table
DELETE FROM trip_sheet_header WHERE trip_sheet_id = 6;
DELETE FROM trip_sheet_header WHERE trip_sheet_id = 7;
DELETE FROM trip_sheet_header WHERE trip_sheet_id = 8;

-- Delete vendor 3 and related data
-- First delete all related records for trip sheets belonging to vendor 3
DELETE tshlup FROM trip_sheet_header_load_unload_points tshlup 
INNER JOIN trip_sheet_header tsh ON tshlup.trip_sheet_id = tsh.trip_sheet_id 
WHERE tsh.vendor_id = 3;

DELETE mp FROM manage_pod mp 
INNER JOIN trip_sheet_header tsh ON mp.trip_sheet_id = tsh.trip_sheet_id 
WHERE tsh.vendor_id = 3;

DELETE lr FROM lr_receipt_header lr 
INNER JOIN trip_sheet_header tsh ON lr.trip_sheet_id = tsh.trip_sheet_id 
WHERE tsh.vendor_id = 3;

-- Now delete trip sheets for vendor 3
DELETE FROM trip_sheet_header WHERE vendor_id = 3;

-- Delete vehicles for vendor 3
DELETE FROM vehicles WHERE vendor_id = 3;

-- IMPORTANT: Delete vendor contact info FIRST (before vendor)
-- This must be done before deleting the vendor due to foreign key constraint
DELETE FROM vendor_contact_info WHERE vendor_id = 3;

-- Delete declaration year info for vendor 3
DELETE FROM declaration_year_info WHERE vendor_id = 3;

-- Ensure foreign key checks are still disabled before deleting vendor
SET FOREIGN_KEY_CHECKS = 0;

-- Finally, delete vendor 3
DELETE FROM vendors WHERE vendor_id = 3;

SET FOREIGN_KEY_CHECKS = 1;

-- Verify deletion
SELECT 
    (SELECT COUNT(*) FROM trip_sheet_header WHERE trip_sheet_id IN (6, 7, 8)) as remaining_trip_sheets_6_7_8,
    (SELECT COUNT(*) FROM trip_sheet_header WHERE vendor_id = 3) as remaining_trip_sheets_vendor_3,
    (SELECT COUNT(*) FROM vendors WHERE vendor_id = 3) as remaining_vendor_3,
    (SELECT COUNT(*) FROM vehicles WHERE vendor_id = 3) as remaining_vehicles_vendor_3,
    (SELECT COUNT(*) FROM vendor_contact_info WHERE vendor_id = 3) as remaining_contact_info_vendor_3,
    (SELECT COUNT(*) FROM declaration_year_info WHERE vendor_id = 3) as remaining_dec_info_vendor_3;

SELECT 'Trip sheets 6, 7, 8 and vendor 3 deleted successfully!' as message;

