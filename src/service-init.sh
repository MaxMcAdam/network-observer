
curl -X PUT http://$(DB_ADMIN_USERNAME):$(DB_ADMIN_PW)@$(DB_URL):5984/_users
curl -X PUT http://$(DB_ADMIN_USERNAME):$(DB_ADMIN_PW)@$(DB_URL):5984/_replicator
curl -X PUT http://$(DB_ADMIN_USERNAME):$(DB_ADMIN_PW)@$(DB_URL):5984/_global_changes

curl -X PUT http://$(DB_ADMIN_USERNAME):$(DB_ADMIN_PW)@$(DB_URL):5984/live-hosts
curl -X PUT http://$(DB_ADMIN_USERNAME):$(DB_ADMIN_PW)@$(DB_URL):5984/auth-hosts
