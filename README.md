# TP0: Docker + Comunicaciones + Concurrencia

### Ejercicio N°1:
Además, definir un script de bash `generar-compose.sh` que permita crear una definición de DockerCompose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```


### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).



### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se puede exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `



### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).



## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base al código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.



### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).



### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

El servidor, por otro lado, deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.



### Ejercicio N°7:
Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo, no podrá responder consultas por la lista de ganadores.
Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.



### Resolucion y Explicacion de Protocolo:
Utilice un procolo de texto, donde cada campo del mensaje esta delimitado por una secuencia de caracteres especiales, cualquier mensaje del cliente tiene la estructura general de:
`<message-type>\n<payload>\n\n` donde `<message-type>` puede tomar los valores **bet**, **ready**, **results**, el contenido de `<payload>` va a depender del `<message-type>`.

Si `<message-type>` es de tipo **bet**, entonces el payload estara formado por un conjunto de apuestas, cada apuesta se serializa de la siguiente forma
`agencyId|firstName|lastName|document|birthDate|number\n`, por ejemplo el mensaje que envia el cliente de la agencia 1 con dos apuestas tendria la siguiente forma

```
bet\n
1|Tomas|Castellano|40770898|13-04-1998|5000\n
1|Ines|Castellano|40770897|14-04-2000|4500\n
\n\n
```

Si `<message-type>` es de tipo **ready** o **results**, entonces el payload estara formado por el id de la agencia, por ejemplo para indicar que todas las apuestas de la agencia fueron
enviadas:
```
ready\n
1\n
\n\n

```
Para preguntas por los resultados finales del sorteo:
```
results\n
1\n
\n\n
```
En todos los casos al final del mensaje estan los caracteres `\n\n` que indican el fin del mensaje, que fue lo que me ayudo a evitar short-reads o short-writes en caso de algun problema en la red.

Por otro lado, el servidor envia mensajes con la forma `<status>,<text>\n` donde `<status>` indica si la operacion fue realizada correctamente y `<text>` es utilizado para enviarle algun tipo de 
mensaje al cliente.

A medida que va guardando los batchs de informacion, les responde:
```
OK,action: receive_message | result: success\n

```
Luego, cuando todas las agencias enviaron sus apuestas y comienzan a preguntar por resultados, el servidor les responde una por una con el siguiente mensaje:

```
OK,action: consulta_ganadores | result: success | cant_ganadores: ${CANT}\n
```
Indicandole que en su agencia ganaron "CANT" apuestas


## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

En caso de que el alumno implemente el servidor Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

### Resolucion y Explicacion de Mecanismos de Sincronización:
Para que el servidor pueda atender varios clientes en paralelo, utilice la libreria de python de multiprocessing para que cada cliente este en un proceso separado y asi evitar utilizar multithreading.
Los mecanismos de sincronizacion que tuve que utilizar son `lock` y `barrier` el primero me permite cuidar la seccion critica que en este caso es el archivo que se va generando en el servidor con todas
las apuestas, solo un proceso por vez puede acceder a esa seccion critica. Luego la barrera me permitio sincronizar los clientes para que esperaran el resultado final del sorteo, la idea es que
el clienteA se quede esperando a que el servidor le comunique el ganador de la loteria, para ello el procesoA (correspondiente al clienteA en el servidor) debe esperar a que todos clientes terminen de enviar sus apuestas, recien
ahi es cuando la barrera deja continuar a los procesos y asi el servidor envia los resultados a cada uno.


## Consideraciones Generales
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios.El _fork_ deberá contar con una sección de README que indique como ejecutar cada ejercicio.

La Parte 2 requiere una sección donde se explique el protocolo de comunicación implementado.
La Parte 3 requiere una sección que expliquen los mecanismos de sincronización utilizados.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deben existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.

(hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Puden obtener un listado del último commit de cada rama ejecutando `git ls-remote`.

Finalmente, se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección provistos [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).
