DMIN=admin
PW=p4ssw0rd
IP=127.0.0.1
PORT=5984

db-setup=$(curl -X GET "http://$IP:$PORT/_all_dbs")
curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/_users"
curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/_replicator"
curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/_global_changes"

curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/live-hosts"

curl -X "POST http://$ADMIN:$PW@$IP:$PORT/mycooldb -H "Content-Type: application/json" -d '{ "somekey": "some value", "anotherkey": "another value" }'"

touch new-scan.txt last-scan.txt delta.txt devices-added.txt devices-removed.txt

./service.sh
