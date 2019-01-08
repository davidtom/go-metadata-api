# Upbound App Metadata API
Golang RESTful API server for storing and searching application metadata expressed via YAML

## Installation
```
go get "github.com/davidtom/go-metadata-api"
```

## Usage
Start the server:
```
make run
```

A port may also be specified if needed:
```
make run PORT=:9000
```

#### Upload Metadata
Upload YAML metadata by making a POST request to `http://localhost:8080/v1/metadata` with a `Content-Type: application/x-yaml` header.

Valid metadata will return a `204` status code:
```
curl --header "Content-Type: application/x-yaml" \
  --request POST \
  --data \
'title: Valid App 1
version: 0.0.1
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
os:
 - linux
 - darwin
metadata:
  label: myapp
description: |
 ### Interesting Title
 Some application content, and description' \
  http://localhost:8080/v1/metadata
```

**Note:** metadata is stored by title and version, so making the request above twice will not store another piece of metadata. To do so, either change the title, increment the version, or both.  

Invalid YAML metadata will be rejected with a `415` status code and a short reason why, if possible.

The metadata below will be rejected because the maintainer email is invalid:
```
curl --header "Content-Type: application/x-yaml" \
  --request POST \
  --data \
'title: App w/ Invalid maintainer email
version: 1.0.1
maintainers:
- name: Firstname Lastname
  email: apptwohotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
 ### blob of markdown
 More markdown' \
  http://localhost:8080/v1/metadata
```

#### Search Metadata
Once stored, metadata can be searched and retrieved by making a GET request to `http://localhost:8080/v1/metadata/search` and specifying search terms as query parameters.

For example, to search for metadata with `version == 0.0.1`, make the following request:
```
curl http://localhost:8080/v1/metadata/search?version=0.0.1
```

In order to search nested data, simply specify the order of keys to access in the field of the query string. Data nested inside arrays is searched independent of index. For example, to search for metadata with a maintainer email of secondmaintainer@gmail.com, the reqest/query parameters are as follows:
```
curl http://localhost:8080/v1/metadata/search?maintainers,email=secondmaintainer@gmail.com
```

**NOTE:** query parameters cannot be duplicated in the requests - and only the last instance of a query parameter will be used

## Bonus
I had so much fun trying to implement search that I wanted to try to cover additional types of nested data aside from just a list of dictionaries (ie maintainers). So, I added two additional - but optional - fields to the metadata:
```yaml
metadata:
  label: mylabel
os:
  - linux
  - darwin
```

To search for stored metadata with the above fields, a request like the following can be made:
```
curl http://localhost:8080/v1/metadata/search?metadata,label=mylabel&os=darwin
```
