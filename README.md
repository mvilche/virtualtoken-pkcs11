## VIRTUAL ETOKEN PKCS11 IMPLEMENTACION - TESTING Y DESARROLLO

### COMPATIBLE CON FEDORA, CENTOS, DEBIAN Y UBUNTU



### INICIAR

#### IMPORTANTE. NO INICIAR LA APLICACION CON USUARIO ROOT, EL SISTEMA SOLICITARÁ LAS CREDENCIALES CON SUDO CUANDO SEA NECESARIO.

La apliucacion detecta de forma automatica si es el primer inicio e instala las dependencias necesarias.

./virtualtoken -start

### DETENER

./virtualtoken -stop


### INICIALIZAR TOKEN Y CERTIFICADO AUTOFIRMADO

#### IMPORTANTE: ANTES DE INICIALIZAR EL TOKEN VERIFIQUE QUE EL USUARIO QUE EJECUTA EL VITUAL TOKEN SEA PARTE DEL GRUPO pkcs11
####

./virtualtoken -init


### DATOS DEL ETOKEN

PIN 123456
SLOT 3

Verificar token:

pkcs11-tool --module PATHTO/softoken/lib/opencryptoki/libopencryptoki.so -L --slot 3 --login --pin 123456