select t.title as title,
       t.done  as done,
       t.date as date
from todos t
where t.id = $1