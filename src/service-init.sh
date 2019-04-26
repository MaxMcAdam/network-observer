PORT=5985

curl -X PUT http://admin:p4ssw0rd@127.0.0.1:$PORT/_users
curl -X PUT http://admin:p4ssw0rd@127.0.0.1:$PORT/_replicator
curl -X PUT http://admin:p4ssw0rd@127.0.0.1:$PORT/_global_changes

curl -X PUT http://admin:p4ssw0rd@127.0.0.1:$PORT/live-hosts
curl -X PUT http://admin:p4ssw0rd@127.0.0.1:$PORT/auth-hosts
