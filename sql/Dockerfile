FROM jbergknoff/postgresql-client

COPY schema.sql .

ENTRYPOINT ["psql"]

CMD ["-1f", "schema.sql"]
