# capital-go

## Summary
This project aims to management some stock assets via tools written by golang.  
This is a developing project, so please take it easy and keep an eye on it.


## Features
| Source | Avaiable Actions | Required Environment Values |
| :--- | :---- | :--- |
| BitFlyer | <ul><li>Show Board</li><li>Show Market</li><li>Show Your Balance (Required Private Key)</li></ul>  | <ul><li>`BITFLYER_API_KEY`</li><li>`BITFLYER_API_SECRET`</li></ul> |

## Build
```
$ make build
```
After building, an executable file `capital-go` will be generated.

## Run
```
# Call BitFlyer API
$ ./capital-go bitflyer markets
```