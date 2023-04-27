# algo-scripts

Scripts for various tasks on algorand.

## Config

Copy `config.example.json` to `config.json` and fill in the values.

## Run script

Edit the values in the script file then run the script with:

```shell
go run <script name>/main.go
```

where `script name` is one of the following:

* cid-to-reserve
* create-asa
* decode-transaction
* get-account-info
* get-asa-ids
* mass-opt-in
* mass-send
