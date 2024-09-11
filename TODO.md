# Why are some activites not synced?

1663 rides are not synced. Maybe add a background job to resync these?

```postgresql
SELECT 
  count(*) 
FROM 
  activity_summary
  LEFT JOIN activity_detail
  ON activity_summary.id = activity_detail.id
WHERE 
  activity_detail.id IS null
  AND activity_type = 'Ride'
LIMIT 100; 
```

# Some rides have 0 segments

If a ride is synced with 0 segments, we should probably resync it. Sometimes
it is correct with 0, but not always. Maybe have a column for the number of times
an activity was synced? So we can just resync some X times.

17674 activities exist

```postgresql
SELECT 
  count(*)
FROM 
  activity_detail
WHERE 
  NOT EXISTS (
    SELECT activities_id FROM segment_efforts
    WHERE segment_efforts.activities_id = activity_detail.id
  );
```
