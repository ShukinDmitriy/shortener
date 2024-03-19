create table if not exists application.public.url
(
    short_key    varchar not null
        constraint url_pk
            primary key
        constraint url_pk_2
            unique,
    original_url varchar
);