# Viltasks

Gestión centralizada de crons


- Sintaxis crons y schedulers
- Envio de notificaciones por email
- Panel de estados
- Api


### Iniciar:

  docker-compose up

### Acceso  

http://localhost:9000


### Sintaxis

Verifique http://localhost:9000/task/sintaxis


### Api rest

Crons en ejecución: http://localhost:9000/api/success

Crons con ejecuciones fallidas: http://localhost:9000/api/failed


### Configuración de email

Configure en conf/app.conf

```shell
[email]
mail.host = smtp.gmail.com
mail.port = 587
mail.user = user@mail.com
mail.password = 12345
mail.disable.tls = false
```

### Configurar base de datos

Configure en conf/app.conf

Por defecto los datos se persisten en database/viltasks.db

```shell
[database]
database.url = database/viltask.db

```

### Configurar Autenticación login

```shell
[auth]
auth.username = admin
auth.password = admin123
```


### Compilar

1 - Clonar

2 - go install -v github.com/revel/cmd/revel@latest

3 - go install -v github.com/revel/revel@latest

4 - go get -v github.com/revel/modules/static

5 - revel run
