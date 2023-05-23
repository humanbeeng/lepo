-- name: GetAuthor :one
select *
from authors
where id = ?
limit 1
;


-- name: ListAuthors :many
select *
from authors
order by name
;

-- name: CreateAuthor :execresult
insert into authors (
  name, bio
) values (?, ?);


-- name: DeleteAuthor :exec
delete from authors
where id = ?
;
