#!/bin/bash

curl  -v -d '{"hello":"world!", "b":true, "c":100}' localhost:8080
#curl  -v -d @../accountancy-init.json localhost:8080