#!/bin/bash

systemctl stop go-ecommerce-api
git pull
go build
systemctl start go-ecommerce-api