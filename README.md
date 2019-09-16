# docker-repo-watcher
[![Build Status](https://travis-ci.org/freundallein/docker-repo-watcher.svg?branch=master)](https://travis-ci.org/freundallein/docker-repo-watcher)

Docker local repository watcher
* Can prune all docker instances
* Can watch repository and delete stale images
* Can update self docker container with new image
* Can remove stale images from docker registry on the same host

## Installation

* ```docker pull freundallein/drwatcher:latest```
* ```docker run -d -v /var/run/docker.sock:/var/run/docker.sock freundallein/drwatcher:latest```

## Work
Starts docker prune crontb jobs with CRONTAB period.  
Start custom image cleaning job,  
for example your app image is ```192.168.1.50:5000/custom-app```
and you want to store only 5 images from  
```
192.168.1.50:5000/custom-app:2019-01-01.1
192.168.1.50:5000/custom-app:2019-01-01.2
192.168.1.50:5000/custom-app:2019-01-01.3
192.168.1.50:5000/custom-app:2019-01-01.4
192.168.1.50:5000/custom-app:2019-01-01.5
192.168.1.50:5000/custom-app:2019-01-01.6
192.168.1.50:5000/custom-app:2019-01-01.latest
```
So, drwatcher will check your local repository and remove stale images.  
After one iteration with IMAGE_AMOUNT=5 your repository would look like this  
```
192.168.1.50:5000/custom-app:2019-01-01.3
192.168.1.50:5000/custom-app:2019-01-01.4
192.168.1.50:5000/custom-app:2019-01-01.5
192.168.1.50:5000/custom-app:2019-01-01.6
192.168.1.50:5000/custom-app:2019-01-01.latest
```

Also, if you want to clean images form docker registry (on the same host),  
you should pass ```-v /your/regitry/path:$REGISTRY_PATH``` and set ```CLEAN_REGISTRY=1```.  
Drwather will discover ```/_manifests/tags/``` and ```/_manifests/revisions/```,  
decide what revisions and tags should be deleted, will delete it,  
then will call registry garbage collect.

## Options
* ```REGISTRY_IP=192.168.1.50```
* ```REGISTRY_PORT=5000```
* ```APP_PREFIX="custom-app"``` (prefix of custom application)
* ```CRONTAB="* * * * * *"``` (starts with seconds)
* ```LOG_LEVEL=DEBUG``` (or ERROR)
* ```PERIOD=60``` (Period of custom image cleaning in seconds)
* ```IMAGE_AMOUNT=5``` (amount of custom images to stay)
* ```AUTOUPDATE=1``` (if you want autoupdate drwatcher)
* ```CLEAN_REGISTRY=1``` (if you want to clean your registry)
* ```REGISTRY_PATH=/var/lib/registry```