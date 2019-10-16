create schema if not exists hello;
create table if not exists hello.birthday(name text primary key, at date);

create or replace function hello.store_birthday(name text, at date) returns void as
$$
insert into hello.birthday(name, at)
values($1, $2)
on conflict (name) do update set at = excluded.at;
$$ language sql strict;


create or replace function hello.retrieve_birthday_message(p_name text, p_today date) returns setof text as $$
with date_offset as (
    select
            at - date_trunc('year', at)::date as birthday,
            p_today - date_trunc('year', p_today)::date as today,
            (date_trunc('year', p_today) + interval '1 year')::date as coming_new_year
    from
        hello.birthday
    where
            name = p_name
),
     days_until_birthday as (
         select
             case when birthday = today
                      then 0
                  when birthday > today
                      then birthday - today
                  else coming_new_year + birthday - p_today
                 end as days
         from date_offset)
select
    case when days = 0
             then format('Happy birthday, %s!', p_name)
         else format('Hello, %s! Your birthday is in %s %s',
                     p_name, days, case when days = 1 then 'day' else 'days' end)
        end as message
from days_until_birthday;
$$ language sql stable strict;
