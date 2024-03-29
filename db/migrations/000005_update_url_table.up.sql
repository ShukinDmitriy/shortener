begin;
alter table public.url
    add if not exists is_deleted bool default false;
commit;