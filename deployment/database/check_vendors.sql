-- Check what vendors exist in the database
SELECT vendor_id, vendor_code, vendor_name, created_at 
FROM vendors 
ORDER BY vendor_code DESC 
LIMIT 10;

-- Count total vendors
SELECT COUNT(*) as total_vendors FROM vendors;

-- Check the highest vendor code
SELECT MAX(vendor_code) as highest_vendor_code FROM vendors;

-- Check if there are any vendors with codes
SELECT vendor_code, COUNT(*) 
FROM vendors 
GROUP BY vendor_code 
ORDER BY vendor_code DESC;


