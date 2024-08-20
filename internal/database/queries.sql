-- name: EndpointsList :many
select * from endpoints;

-- name: EndpointsCreate :exec
insert into endpoints (url, expected_status) values (?, ?);

-- name: ChecksForEndpoint :many
select * from checks where endpoint_id = ? order by created_at desc limit ?;

-- name: ChecksCreate :exec
insert into checks (endpoint_id, status, response_time) values (?, ?, ?);

-- name: EndpointsDelete :exec
delete from endpoints where id = ?;

-- name: ChecksDelete :exec
delete from checks where endpoint_id = ?;

-- name: ChecksCleanup :exec
delete from checks where created_at < (strftime('%s', 'now', ?));

-- name: Changes :one
select changes();
