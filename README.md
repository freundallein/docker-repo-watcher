# docker-repo-watcher
[![Build Status](https://travis-ci.org/freundallein/docker-repo-watcher.svg?branch=master)](https://travis-ci.org/freundallein/docker-repo-watcher)

Docker local repository watcher
* Can prune all docker instances
* Can watch repository and delete stale images
* Can update self docker container with new image

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
After one iteration with IMAGE_AMOUNT=5 your repositroy will be looked like this  
```
192.168.1.50:5000/custom-app:2019-01-01.3
192.168.1.50:5000/custom-app:2019-01-01.4
192.168.1.50:5000/custom-app:2019-01-01.5
192.168.1.50:5000/custom-app:2019-01-01.6
192.168.1.50:5000/custom-app:2019-01-01.latest
```

## Options
* ```REGISTRY_IP=192.168.1.50```
* ```REGISTRY_PORT=5000```
* ```APP_REFIX="custom-app"``` (prefix of custom application)
* ```CRONTAB="* * * * * *"``` (starts with seconds)
* ```LOG_LEVEL=DEBUG``` (or ERROR)
* ```PERIOD=60``` (Period of custom image cleaning in seconds)
* ```IMAGE_AMOUNT=5``` (amount of custom images to stay)
* ```AUTOUPDATE=1``` (if you want autoupdate drwatcher)