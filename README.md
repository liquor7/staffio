# staffio
## define

    An OAuth2 server with management for enterprise employees


## prepare

### database
````sh
docker run -e DB_NAME=staffio -e DB_USER=staffio -e DB_PASS=mypassword -d --name staffio-db lcgc/postgresql:9.5.2
cat database/schema.sql | docker exec -i staffio-db staffio staffio
cat database/init.sql | docker exec -i staffio-db staffio staffio

-- demo client
echo "INSERT INTO oauth_client VALUES(1, '1234', 'Demo', 'aabbccdd', 'http://localhost:3000/appauth', '{}', now());" | docker exec -i staffio-db staffio staffio


````

### checkout

````sh
mkdir -p $GOPATH/src/lcgc/platform
cd $GOPATH/src/lcgc/platform
git clone https://github.com/liut/keeper.git
git clone https://github.com/liut/staffio.git
````


## launch development

````sh
go get -u github.com/ddollar/forego
go get -u github.com/ddollar/rerun

forego start
````
