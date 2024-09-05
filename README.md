# TP0: Docker + Comunicaciones + Concurrencia

### Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

En caso de que el alumno implemente el servidor Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).


#### Ejecucion
La carptea data/datasets.zip dentro de client debe existir

Luego ejecutar
```
make docker-compose-up
```

Luego, ejecutar y verificar con
```
make docker-compose-logs
```
