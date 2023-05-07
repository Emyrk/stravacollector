 
 -- All hugel segment efforts.
 SELECT merged.athlete_id,
    merged.segment_ids,
    merged.total_time_seconds,

    -- (SELECT name FROM segments WHERE id = any ARRAY(
    --         (SELECT elem->>'segment_id'
    --   FROM json_array_elements(merged.efforts) elem))
    -- ) AS segnames,


    ARRAY((SELECT elem->>'segment_id'
FROM json_array_elements(merged.efforts) elem)) as segments_done,

    ARRAY(SELECT name FROM segments WHERE segments.id :: text = ANY(
ARRAY((SELECT elem->>'segment_id'
FROM json_array_elements(merged.efforts) elem))


    )),
    merged.efforts
   FROM ( SELECT hugel_efforts.athlete_id,
            array_agg(hugel_efforts.segment_id) AS segment_ids,
            sum(hugel_efforts.elapsed_time) AS total_time_seconds,
            json_agg(json_build_object('activity_id', hugel_efforts.activities_id, 'effort_id', hugel_efforts.id, 'start_date', hugel_efforts.start_date, 'segment_id', hugel_efforts.segment_id, 'elapsed_time', hugel_efforts.elapsed_time, 'moving_time', hugel_efforts.moving_time, 'device_watts', hugel_efforts.device_watts, 'average_watts', hugel_efforts.average_watts)) AS efforts
           FROM ( SELECT DISTINCT ON (segment_efforts.athlete_id, segment_efforts.segment_id) segment_efforts.id,
                    segment_efforts.athlete_id,
                    segment_efforts.segment_id,
                    segment_efforts.name,
                    segment_efforts.elapsed_time,
                    segment_efforts.moving_time,
                    segment_efforts.start_date,
                    segment_efforts.start_date_local,
                    segment_efforts.distance,
                    segment_efforts.start_index,
                    segment_efforts.end_index,
                    segment_efforts.device_watts,
                    segment_efforts.average_watts,
                    segment_efforts.kom_rank,
                    segment_efforts.pr_rank,
                    segment_efforts.updated_at,
                    segment_efforts.activities_id
                   FROM segment_efforts
                  WHERE (segment_efforts.segment_id = ANY (ARRAY( SELECT competitive_routes.segments
                           FROM competitive_routes
                          WHERE (competitive_routes.name = 'das-hugel'::text))))
                  ORDER BY segment_efforts.athlete_id, segment_efforts.segment_id, segment_efforts.elapsed_time) hugel_efforts
          GROUP BY hugel_efforts.athlete_id) merged
  WHERE true or (merged.segment_ids @> ARRAY( SELECT competitive_routes.segments
           FROM competitive_routes
          WHERE (competitive_routes.name = 'das-hugel'::text))
          )