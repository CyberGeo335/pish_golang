curl -i -X POST http://178.72.139.210:8088/notes \
-H 'Content-Type: application/json' \
-d '{"title":"Hello","content":"World"}'

curl -i http://localhost:8080/notes/1
