
alter table application.public.url
    add if not exists correlation_id varchar;
