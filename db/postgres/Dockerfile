# base image
FROM postgres:13-alpine

ARG CACHE_BUST=1

COPY init.sql /docker-entrypoint-initdb.d/
