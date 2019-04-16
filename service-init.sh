DMIN=admin
PW=p4ssw0rd
IP=127.0.0.1
PORT=5984

db-setup=$(curl -X GET "http://$IP:$PORT/_all_dbs")
curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/_users"
curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/_replicator"
curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/_global_changes"

curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/offline-authorized-hosts"
curl -X PUT "http://$ADMIN:$PW@$IP:$PORT/online-hosts"

touch new-scan.txt last-scan.txt delta.txt devices-added.txt devices-removed.txt

while read line; do
  curl -X POST http://127.0.0.1:5984/offline-authorized-hosts -H "Content-Type: application/json" -d '{"device-name":"$line", "last-ip":"null", "time-discovered":null, "authorized":"true", "always-on":"false"}'
done < authorized-devices.txt

./service.sh
