#!/bin/bash

sed -i "s/__UUID__/$(UUID)/g" config.json
./vmess -c config.json

./v2ray -c config.json