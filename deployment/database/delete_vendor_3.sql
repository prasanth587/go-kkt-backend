-- Delete vendor 3 and all related data
-- This script deletes vendor 3 and all associated records

SET FOREIGN_KEY_CHECKS = 0;

-- Delete vendor 3 related data (child tables first)
-- First delete all trip sheet related records for vendor 3
DELETE tshlup FROM trip_sheet_header_load_unload_points tshlup 
INNER JOIN trip_sheet_header tsh ON tshlup.trip_sheet_id = tsh.trip_sheet_id 
WHERE tsh.vendor_id = 3;

DELETE mp FROM manage_pod mp 
INNER JOIN trip_sheet_header tsh ON mp.trip_sheet_id = tsh.trip_sheet_id 
WHERE tsh.vendor_id = 3;

DELETE lr FROM lr_receipt_header lr 
INNER JOIN trip_sheet_header tsh ON lr.trip_sheet_id = tsh.trip_sheet_id 
WHERE tsh.vendor_id = 3;

-- Delete trip sheets for vendor 3
DELETE FROM trip_sheet_header WHERE vendor_id = 3;

-- Delete vehicles for vendor 3
DELETE FROM vehicles WHERE vendor_id = 3;

-- Delete vendor contact info for vendor 3
DELETE FROM vendor_contact_info WHERE vendor_id = 3;

-- Delete declaration year info for vendor 3
DELETE FROM declaration_year_info WHERE vendor_id = 3;

-- Finally, delete vendor 3
DELETE FROM vendors WHERE vendor_id = 3;

SET FOREIGN_KEY_CHECKS = 1;

-- Verify deletion
SELECT 
    (SELECT COUNT(*) FROM vendors WHERE vendor_id = 3) as remaining_vendor_3,
    (SELECT COUNT(*) FROM trip_sheet_header WHERE vendor_id = 3) as remaining_trip_sheets,
    (SELECT COUNT(*) FROM vehicles WHERE vendor_id = 3) as remaining_vehicles,
    (SELECT COUNT(*) FROM vendor_contact_info WHERE vendor_id = 3) as remaining_contact_info,
    (SELECT COUNT(*) FROM declaration_year_info WHERE vendor_id = 3) as remaining_dec_info;

SELECT 'Vendor 3 deleted successfully!' as message;

