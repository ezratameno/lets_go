FROM mysql:8.0.23
# copy all the files that ends in .sql to the db entrypoint
COPY ./pkg/models/mysql/*.sql /docker-entrypoint-initdb.d/