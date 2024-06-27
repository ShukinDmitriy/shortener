begin;
alter table public.url
    add if not exists user_id varchar;
commit;