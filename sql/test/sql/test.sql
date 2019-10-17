\i schema.sql
TRUNCATE hello.birthday;

SELECT hello.store_birthday('Nguyễn Văn A', '0001-01-01');
SELECT * FROM hello.retrieve_birthday_message('Nguyễn Văn A', current_date);
SELECT hello.store_birthday('Nguyễn Văn A', '2016-02-29');
SELECT * FROM hello.retrieve_birthday_message('Nguyễn Văn A', '2020-02-29');
SELECT * FROM hello.retrieve_birthday_message('Nguyễn Văn A', '2020-03-01');
SELECT * FROM hello.retrieve_birthday_message('Nguyễn Văn A', '2021-02-28');
SELECT * FROM hello.retrieve_birthday_message('Nguyễn Văn A', '2021-03-01');
