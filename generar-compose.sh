#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Uso: $0 <file_name> <clients_number>"
    exit 1
fi

# Asigna los parámetros a variables
file_name=$1
clients_number=$2

echo "Nombre del archivo de salida: $file_name"
echo "Cantidad de clientes: $clients_number"
python3 ./utils/yaml/yaml_generator.py $file_name $clients_number