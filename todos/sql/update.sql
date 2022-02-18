update todos t
set title = :title,
    done  = :done,
    date  = :date
where t.id = :id