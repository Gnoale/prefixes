FROM postgres 

COPY repository/schema.sql /docker-entrypoint-initdb.d/
