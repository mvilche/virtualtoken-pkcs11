# Virtual eToken pkcs11 implementación para desarrollo y testing.

Simula un etoken conectado al equipo para la firma de documentos

# Descarga

Descargue la última versión estable.

```sh
$ mkdir virtualtoken/ && cd virtualtoken/
$ wget https://github.com/mvilche/virtualtoken-pkcs11/releases/download/v1.3/virtualtoken
$ chmod +x virtualtoken
```

# Instalación

## Requisitos previos

- Compatible unicamente con Linux en Ubuntu, Debian, Fedora y CentOS
- Paquetes sudo y pidof instalados en el equipo antes de comenzar la instalación.
- El usuario que ejecutará la instalación debe pertenecer al grupo sudo para la instalación de dependencias.


## Iniciar

si es el primer inicio el software realizará la instalación de las dependencias necesarias.


```sh
$ ./virtualetoken -start
```

> Finalizada la instalación del virtualtoken 
> se deberá reiniciar la sesión del usuario
> o reiniciar el equipo para que los cambios
> surjan efecto antes de la inicialización
> del virtualtoken.


## Inicializar virtualtoken

El proceso de inicialización definira un etoken virtual en el slot 3 con un pin númerico (123456).
Adicionalmente creara un juego de llaves y certificados autofirmados que serán importados dentro del etoken.

```sh
$ ./virtualetoken -start
$ ./virtualetoken -init
```

## Detener

Detiene el virtualtoken

```sh
$ ./virtualetoken -stop
```

# Datos del virtualtoken

* Pin: 123456
* Slot: 3
* Driver: PATH/TO/softoken/lib/opencryptoki/libopencryptoki.so


# Comandos utiles pkcs11-tools

## Acceder y listar etoken

```sh
pkcs11-tool --module PATH/TO/softoken/lib/opencryptoki/libopencryptoki.so -L --slot 3 --login --pin 123456
```
## Ver certificados dentro del virtualtoken

```sh
pkcs11-tool --module PATH/TO/softoken/lib/opencryptoki/libopencryptoki.so -L --slot 3 --list-objects --type cert
```
## Generar nuevo certificado e importarlo dentro del virtualetoken

### Generar clave privada y certificado

```sh
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout private.key -out cert.pem -subj "/CN=testing"
```
### Convertir certificado a DER

```sh
openssl x509 -inform pem -outform der -in cert.pem -out cert.der
```
### Convertir clave privada a DER

```sh
openssl rsa -inform pem -outform der -in private.key -out key.der
```
### Importar clave privada DER en virtualtoken

```sh
pkcs11-tool --module PATH/TO/softoken/lib/opencryptoki/libopencryptoki.so --slot 3 --login --pin 123456 --write-object key.der --type privkey --id 20 --label 'mi_cert'
```
### Importar certificado DER en virtualtoken

```sh
pkcs11-tool --module PATH/TO/softoken/lib/opencryptoki/libopencryptoki.so --slot 3 --login --pin 123456 --write-object cert.der --type cert --id 20 --label 'mi_cert'
```
### Verificar nuevo certificado

```sh
pkcs11-tool --module PATH/TO/softoken/lib/opencryptoki/libopencryptoki.so -L --slot 3 --list-objects --type cert

```
