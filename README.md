# Breeze

[![Go Report Card](https://goreportcard.com/badge/github.com/GaruGaru/breeze)](https://goreportcard.com/report/github.com/GaruGaru/breeze)
[![Travis Card](https://travis-ci.org/GaruGaru/breeze.svg?branch=master)](https://travis-ci.org/GaruGaru/breeze)
[![MicroBadger Size](https://img.shields.io/microbadger/image-size/garugaru/breeze)](https://cloud.docker.com/u/garugaru/repository/docker/garugaru/breeze)
 

Smart, Temperature based fan controller for raspberry pi using gpio written in go

## Schema for fan control using transistor 
[Schema (using gpio 12)](https://4.bp.blogspot.com/-iiFMc_3jrdU/XGDhy20-EbI/AAAAAAAAAlc/wJ_u3oXGZqkkS8kh9VBjhQdyrmaXKpRdQCK4BGAYYCw/s1600/gpio-fan.png)

## Usage

example: 

When the detected temperature reaches the *target-temperature* (60°) the gpio pin n.12 will be set to HIGH until the detected temperature
is 15% lower than the target-temperature 


```bash
breeze controller --gpio-pin=12 --target-temperature=60 --temperature-cooldown-percent=15
```

## Deploy 

### Docker
```bash
docker run --privileged -v /dev/gpiomem:/dev/gpiomem -v /sys/class/gpio:/sys/class/gpio garugaru/breeze:arm-latest controller 
```

### kubernetes
```bash
kubectl apply -f https://raw.githubusercontent.com/GaruGaru/breeze/master/kubernetes/0_namespace.yml
kubectl apply -f https://raw.githubusercontent.com/GaruGaru/breeze/master/kubernetes/1_deployment.yml
```



 