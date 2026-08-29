FROM migrate/migrate:v4.18.3

COPY db/migrations /migrations

ENTRYPOINT ["migrate"]
CMD ["-path", "/migrations", "-database", "postgres://postgres:password@pgtestdb:5432/postgres?sslmode=disable", "up"]
