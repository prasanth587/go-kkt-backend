-- Check if vendors table is truly empty
SELECT COUNT(*) as total_vendors FROM vendors;

-- Check if there are any vendor codes in the database
SELECT vendor_code FROM vendors ORDER BY vendor_code DESC LIMIT 10;

-- Check trip_sheet_num again
SELECT * FROM trip_sheet_num WHERE year = 'VEN';

-- If vendors exist, show the highest vendor code
SELECT MAX(vendor_code) as highest_vendor_code FROM vendors;


