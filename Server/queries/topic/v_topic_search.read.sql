-- 搜索框模糊查询
select *
from v_topic
where title LIKE '%{{.match}}%'
   OR content LIKE '%{{.match}}%'
   OR nickname LIKE '%{{.match}}%'
   OR array_to_string(label, ',') LIKE '%{{.match}}%'
ORDER BY updated_at DESC
LIMIT 30
