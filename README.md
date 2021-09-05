# reload-config
reload config and shutdown app using http request  
wait for subscriber to finish processing  
synchronize all go routine before shutdown  

### rabbitmq
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.9-management  