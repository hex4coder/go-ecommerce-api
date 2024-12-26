#!/bin/bash

systemctl stop go-ecommerce-api
go build
systemctl start go-ecommerce-api