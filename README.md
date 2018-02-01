# Sproot
Engine for store data in graph database and send notifiy.


```docker
docker pull dgraph/dgraph

# Directory to store data in. This would be passed to `-v` flag.
mkdir -p /tmp/data

# Run Dgraph Zero
docker run -it -p 8080:8080 -p 9080:9080 -v /tmp/data:/dgraph --name diggy dgraph/dgraph dgraph zero --port_offset -2000

# Run Dgraph Server
docker exec -it diggy dgraph server --memory_mb 2048 --zero localhost:5080

```

```
// For test
go get
go test ./...
```

## With DockerCompose for use with docker-compose.yaml
```
docker-compose up -d

docker-compose stop
```

## Example of Instruction for company

```$xslt
{  
   "Data":{  
      "Language":"en",
      "Company":{  
         "ID":"0x2786",
         "Name":"Company test name",
         "IRI":""
      },
      "Category":{  
         "ID":"",
         "Name":""
      },
      "City":{  
         "ID":"0x2788",
         "Name":"Test city"
      },
      "Page":{  
         "uid":"0x2789",
         "path":"smartfony-i-svyaz/smartfony-205",
         "pageInPaginationSelector":".pagination-list .pagination-item",
         "previewImageOfSelector":"",
         "pageParamPath":"/f/page=",
         "cityParamPath":"?cityId=",
         "cityParam":"CityCZ_975",
         "itemSelector":".grid-view .product-tile",
         "nameOfItemSelector":".product-tile-title",
         "cityInCookieKey":"",
         "cityIdForCookie":"",
         "priceOfItemSelector":".product-price-current"
      }
   },
   "Message":"Parse products of company"
}
```