#!/bin/bash

sleep 2 && goose -dir "./migrations" postgres "host=pg-auth port=5432 dbname=auth user=auth-user password=auth-password sslmode=disable" up -v
sleep 2 && goose -dir "./migrations" postgres "host=pg-test-auth port=5432 dbname=auth user=auth-user password=auth-password sslmode=disable" up -v
