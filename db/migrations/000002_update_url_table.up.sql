begin;
alter table public.url
    add if not exists correlation_id varchar;
commit;