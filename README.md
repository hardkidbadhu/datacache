# Datacache
Data Caching Service implementation in go backed by redis, mongodb Kafka

# Checklist

Store data in cache with persistence backup from DB [x] - written as stand alone job that is executed evey 4 minutes in the kubernetes cron resis -> mongodb
Put data in to cache service using REST API [x] - write the data into cache 
Read data from cache service using REST API (support pagination) [x] - Reads the data from cache pagination supported
Reload data from DB with a notification from message broker like Kafka [X] - while reading data from cache if not present in cache it will be fetched from mongodb and the data requested will be cached in redis for subsequent same request.

Documentation using Swagger [x] - documentation availabe navigate to deployments/docs/swagger.yml

The code should be compiled and deployable in minikube [x] - navigate to deployments/k8s/ for deployment files.
