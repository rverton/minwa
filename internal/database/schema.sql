create table if not exists endpoints (
    id integer primary key, -- alias for rowid
    url text not null,
    expected_status integer not null,

    created_at integer not null default (strftime('%s', 'now'))
) strict;

create table if not exists checks (
    endpoint_id integer not null,
    status integer not null,
    response_time integer not null,

    created_at integer not null default (strftime('%s', 'now')),

    foreign key(endpoint_id) references endpoints(id)
) strict;

create index if not exists checks_endpoints_idx on checks (endpoint_id, created_at);
