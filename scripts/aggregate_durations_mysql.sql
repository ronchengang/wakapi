SELECT project, language, editor, operating_system, machine, branch, SUM(diff) as 'sum'
FROM (SELECT project,
             language,
             editor,
             operating_system,
             machine,
             branch,
             TIME_TO_SEC(LEAST(TIMEDIFF(time, LAG(time) over w), '00:00:00')) as 'diff'  -- time constant ~ heartbeats padding (none by default, formerly 2 mins)
      FROM heartbeats
      WHERE user_id = 'n1try'
      WINDOW w AS (ORDER BY time)) s2
WHERE diff IS NOT NULL
GROUP BY project, language, editor, operating_system, machine, branch;

-- for mysql 5.7
-- SELECT project, language, editor, operating_system, machine, branch, SUM(GREATEST(1, diff)) as 'sum'
-- FROM (
--          SELECT
--              h1.project,
--              h1.language,
--              h1.editor,
--              h1.operating_system,
--              h1.machine,
--              h1.branch,
--              TIME_TO_SEC(LEAST(TIMEDIFF(h1.time,
--                                         (SELECT h2.time
--                                          FROM heartbeats h2
--                                          WHERE h2.time < h1.time
--                                            AND h2.user_id = h1.user_id
--                                          ORDER BY h2.time DESC
--                                          LIMIT 1)
--                                ), '00:02:00')) as 'diff'
--          FROM heartbeats h1
--          WHERE h1.user_id = 'ron'
--      ) s2
-- WHERE diff IS NOT NULL
-- GROUP BY project, language, editor, operating_system, machine, branch;