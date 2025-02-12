#!/bin/bash


cp .env ../
systemctl stop go-ecommerce-api
git stash
git pull
go build
cp ../.env .
systemctl start go-ecommerce-api
systemctl status go-ecommerce-api