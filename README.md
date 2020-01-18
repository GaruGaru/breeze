# Breeze

Smart, Temperature based fan controller for raspberry pi using gpio written in go

## Schema for fan control using transistor 
[Schema (using gpio 12)](https://4.bp.blogspot.com/-iiFMc_3jrdU/XGDhy20-EbI/AAAAAAAAAlc/wJ_u3oXGZqkkS8kh9VBjhQdyrmaXKpRdQCK4BGAYYCw/s1600/gpio-fan.png)

## Usage

```bash
breeze controller --gpio-pin=12 --target-temperature=60
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



 