# capital-go

## Summary
This project aims to management some stock assets via tools written by golang.  
This is a developing project, so please take it easy and keep an eye on it.


## Features
### BitFlyer
If you want to use `Authorization Required` actions, you must need to set these values to Environment Variables.
- `BITFLYER_API_KEY`
- `BITFLYER_API_SECRET`

| Actions | Authorization |
| :---- | :--- |
| Show Board | - |
| Show Market| - |
| Show Balance | Required |
| Send Order | Required |

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
