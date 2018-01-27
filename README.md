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
[  
   {  
      "uid":"0x1eb",
      "instructionLanguage":"en",
      "instructionIsActive":true,
      "has_pages":[  
         {  
            "uid":"0x1ed",
            "path":"smartfony-i-svyaz/smartfony-205",
            "pageInPaginationSelector":".pagination-list .pagination-item",
            "pageParamPath":"/f/page=",
            "cityParamPath":"?cityId=",
            "cityParam":"CityCZ_975",
            "itemSelector":".grid-view .product-tile",
            "nameOfItemSelector":".product-tile-title",
            "priceOfItemSelector":".product-price-current"
         }
      ],
      "has_city":[  
         {  
            "uid":"0x1ec",
            "cityName":"Test city",
            "cityIsActive":true
         }
      ],
      "has_company":[  
         {  
            "uid":"0x1ea",
            "companyIri":"",
            "companyName":"Test company",
            "has_category":null,
            "companyIsActive":true
         }
      ],
      "has_category":[  
         {  
            "uid":"0x1c9",
            "categoryName":"Test category",
            "categoryIsActive":true,
            "belongs_to_company":null,
            "has_product":null
         }
      ]
   }
]
```