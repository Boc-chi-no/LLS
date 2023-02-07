# LinkShortener
### A URL shortener written in GO

~~This was built using Docker-Compose. In order to run in your machine, just clone the repository and run:~~
~~* sudo docker-compose up --build (first time, then you go up without the --build, which is much faster)~~

~~If you don't have docker installed in your computer, please download and install it at:~~
~~https://www.docker.com/~~

## Instructions
### Captcha

To perform a create/manage operation you need to create Captcha first, just http GET to http://localhost:8040/api/captcha, The API will return the following:
```json
{
  "code":0,
  "data":{
    "pic":"data:image/png;base64,....."
  },
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```

Then, a Base64-encoded challenge image and a cookie identifying the Session are returned

### Generate

To shorten a URL, just http POST to http://localhost:8040/api/generate_link with the following json payload (example):

```json
{
  "link":"http://127.0.0.1:8040/", //Original URL
  "captcha":"8" //Captcha answer
}
```

The api will return the following:

```json
{
  "code":0,
  "data":{
    "hash":"18nfqL", //shortened URL Hash
    "token":"IKmXKMrVtBOvdibt" //Manage Password
  },
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```

The token is your subsequent credentials for managing the link, and the hash is the shortened URL Hash

### Redirect

just http GET to http://localhost:8040/s/:hash. and you will get the URL redirection (example):
```
http://localhost:8040/s/18nfqL
```


### Statistics

The application will record the visit and write it to the database, just http POST to http://localhost:8040/api/stats_link with the following json payload (example):

```json
{
  "hash": "18nfqL", //shortened URL Hash
  "token": "IKmXKMrVtBOvdibt", //Manage Password
  "captcha": "25" //Captcha answer
  "page": 1, // Page number of current visit(A positive integer)
  "size": 50 //Size per page(integers from 1-100)
}
```
The api will return the following:

```json
{
  "code":0,
  "data":{
    "current":1, //current page
    "size":50, //Size set by request
    "pages":1, //Total number of pages
    "total":1, //Total number of results
    "records":[
      {
        "Hash":"18nfqL", //HASH of the query
        "IP":"127.0.0.1", //IP of the visitor
        "Header":{ //Request Header of the visitor
          "Accept-Encoding":[
            "gzip, deflate, br"
          ]
        },
        "Country":"Local Address", //The country indicated by the visitor's IP
        "Area":"Local Address", //The area indicated by the visitor's IP
        "Browser":"Chrome", //The Browser indicated by the visitor's UA
        "BrowserVersion":"109.0.0",//The Browser Version indicated by the visitor's UA
        "OS":"Windows", //The OS indicated by the visitor's UA
        "OSVersion":"10", //The OS Version indicated by the visitor's UA
        "Device":"Other", //The Device indicated by the visitor's UA
        "Created":1675143659 //Access time (seconds timestamp)
      }
    ]
  },
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```
It will show detailed data about the URL accessed.

### Delete
If the link needs to be removed, just http POST to http://localhost:8040/api/delete_link with the following json payload (example):

```json
{
  "hash": "18nfqL", //shortened URL Hash
  "token": "IKmXKMrVtBOvdibt", //Manage Password
  "captcha": "32" //Captcha answer
}
```
The api will return the following:

```json
{
  "code":0,
  "data":null,
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```
The link will be marked for deletion, but note that it can still be queried for statistics using the administrative password.

