FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY todo_server /app/todo_server
COPY web /app/web

RUN chmod +x /app/todo_server

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/db/scheduler.db
ENV TODO_PASSWORD=my_secret_pass

EXPOSE 7540

CMD ["./todo_server"]