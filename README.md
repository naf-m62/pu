# PU
Custom microservice skeleton based on uber fx

### Used  
- DB - postgres
- Broker - rabbitMQ

### Launch  
app -c .app/config.yaml

### Add handlers
- For http - add new handler and route to cmd/handlers/http_handlers/handler.go
- For RabbitMQ - add new handler, exchange, routing key to cmd/handlers/rmq_handlers/handler.go
- For grpc - add a new handler to proto and implement it in cmd/handlers/grpc_handlers/public/handler.go