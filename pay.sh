#!/bin/bash
ID=$1
PRICE=$2
json_data="{\"id\":\"$ID\",\"price\": $PRICE}"
echo $json_data
curl -X POST -d "$json_data" http://localhost:7777/pay
