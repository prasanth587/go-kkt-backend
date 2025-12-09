-- Check what's in trip_sheet_num table for VENDOR
SELECT * FROM trip_sheet_num WHERE year = 'VEN';

-- Check all entries in trip_sheet_num
SELECT * FROM trip_sheet_num ORDER BY id DESC LIMIT 10;

-- Delete VEN entries from trip_sheet_num (this is what's causing the issue)
DELETE FROM trip_sheet_num WHERE year = 'VEN';

-- Verify it's deleted
SELECT COUNT(*) as remaining_ven_entries FROM trip_sheet_num WHERE year = 'VEN';

SELECT 'trip_sheet_num VEN entries deleted! Now vendor code should start from VEN-000001' as message;


