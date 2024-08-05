# IPSTAT

<!--Inspired by https://ifconfig.me/.-->

## Features
* speed test

## Endpoints
* GET `/health`: returns 200OK
* GET `/` or `/ip`: plaintext with ip
* GET `/ua` or `/useragent`: plaintext with user-agent
* GET `/forwarded`: plaintext with header "X-FORWARDED-FOR"
* GET `/all`: plaintext with all headers
* GET `/all.json`: json with all headers

