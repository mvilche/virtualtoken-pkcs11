package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func main() {

	checkRoot()

	AskSudo()

	fmt.Print("***********************************************************************\n")
	fmt.Print("DESC: VIRTUAL TOKEN - FIRMA DE DOCUMENTOS Y GESTION DE CERTIFICADOS\n")
	fmt.Print("AUTOR: MARTIN FABRIZZIO VILCHE\n")
	fmt.Print("VERSION: 1.3\n")
	fmt.Print("***********************************************************************\n")

	err := createLogFile()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	command := DetectSO()

	ready := ReadyInstall()

	fla, err := getFlag()
	if err != nil {
		fmt.Println("Error ", err.Error())
		ErrorLog.Printf(err.Error())
		os.Exit(1)
	}

	if !ready {
		fmt.Print("--> SOFTWARE NO INSTALADO!\n")
		fmt.Print("--> COMIENZA INSTALCION DEL SOFTWARE\n")
		fmt.Print("--> EL PROGRAMA EJECUTARA COMANDOS CON SUDO PARA INSTALAR LAS DEPENDENCIAS NECESARIAS, INGRESE LA CONTRASENA CUANDO LO SOLICITE.\n")
		err := InstallDependencies(command)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print("--> INSTLACION DE DEPENDENCIAS OK\n")

		err = FixGroupUser()
		if err != nil {
			fmt.Print(err)
		}

		fmt.Print("--> CREACION DE GRUPO Y PERMISOS CORRECTA\n")

		err = DownloadSoftToken()
		if err != nil {
			fmt.Print(err)
		}

		fmt.Print("--> DESCARGA DE SOFTWARE CORRECTA\n")

		err = InstallSoftToken()
		if err != nil {
			fmt.Print(err)
		}

		fmt.Print("--> INSTALACION DE SOFTOKEN CORRECTA\n")

		FixPermis()

		fmt.Print("--> CONFIGURACION DE PERMISOS CORRECTA\n")
		fmt.Print("--> REINICIE EL EQUIPO PARA QUE SU USUARIO PERTENEZCA AL GRUPO pkcs11 ANTES DE INICIALIZAR EL VIRTUAL TOKEN\n")

	} else {

		if fla.Start {

			StartToken()

		}

		if fla.Stop {
			StopToken()
		}

		if fla.Init {
			InitToken()
		}

	}

}

func DetectSO() string {

	var plataform string

	ubuntu := "export DEBIAN_FRONTEND=noninteractive && export DEBCONF_NONINTERACTIVE_SEEN=true && sudo -E apt-get update && sudo -E apt-get install build-essential git openssl opensc autoconf automake make libtool pkg-config m4 flex bison libldap2-dev libssl-dev -y"
	centos7 := "sudo yum install git openssl opensc openssl-devel autoconf automake make libtool pkg-config m4 flex bison openldap-devel -y"
	centos8 := "sudo dnf install git openssl opensc openssl-devel autoconf automake make libtool pkg-config m4 flex bison openldap-devel -y"
	arch := "yay -Syyuu"

	outOS, err := exec.Command("bash", "-c", "awk -F= '/^NAME/{print $2}' /etc/os-release").Output()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	escapeSO := strings.Replace(string(outOS), "\"", "", -1)
	escapeSO = strings.Trim(escapeSO, "")

	outVersion, err := exec.Command("bash", "-c", "awk -F= '/^VERSION_ID/{print $2}' /etc/os-release").Output()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	escapeVersion := strings.Replace(string(outVersion), "\"", "", -1)
	escapeVersion = strings.Trim(escapeVersion, "")

	if escapeSO == "Ubuntu\n" || escapeSO == "Debian\n" {
		plataform = ubuntu
	}

	if escapeSO == "CentOS Linux\n" {
		if escapeVersion == "8\n" {
			plataform = centos8
		}
		if escapeVersion == "7\n" {
			plataform = centos7
		}
	}

	if escapeSO == "Fedora\n" {
		plataform = centos8
	}
	if escapeSO == "Arch Linux\n" {
		fmt.Print("--> PLATAFORMA ENCONTRADA: ", escapeSO)
		fmt.Print("--> PLATAFORMA NO SOPORTADA!\n")
		plataform = arch
		//os.Exit(1)
	}

	_, err1 := exec.LookPath("bash")
	if err1 != nil {
		fmt.Print("--> ERROR BASH NO ENCONTRADO EN PATH.\n")
		os.Exit(1)
	}

	_, err2 := exec.LookPath("sudo")
	if err2 != nil {
		fmt.Print("--> ERROR SUDO NO ENCONTRADO EN PATH.\n")
		os.Exit(1)
	}

	_, err3 := exec.LookPath("pidof")
	if err3 != nil {
		fmt.Print("--> ERROR PIDOF NO ENCONTRADO EN PATH.\n")
		os.Exit(1)
	}

	return plataform
}

func ReadyInstall() bool {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(path + "/softoken"); !os.IsNotExist(err) {
		return true

	}

	return false
}

func InstallDependencies(command string) error {

	cmd := exec.Command("bash", "-c", command)

	// Crear pipe para capturar commando en vivo

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t> %s\n", scanner.Text())
			ExecutionLog.Printf("\t%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()

	return err
}

func DownloadSoftToken() error {

	git := "https://github.com/mvilche/opencryptoki.git"
	cmd := exec.Command("bash", "-c", "git clone "+git+" src")

	// Crear pipe para capturar commando en vivo

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t> %s\n", scanner.Text())
			ExecutionLog.Printf("\t%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return err
}

func FixGroupUser() error {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("bash", "-c", "sudo groupadd pkcs11 && sudo usermod -aG pkcs11 "+user.Username+"")

	// Crear pipe para capturar commando en vivo

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t> %s\n", scanner.Text())
			ExecutionLog.Printf("\t%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return err
}

func FixPermis() {

	path, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	user, err := user.Current()
	if err != nil {
		os.Exit(1)
	}

	fmt.Print("--> CAPTURANDO DATOS DEL USUARIO:\n")
	fmt.Print("--> " + user.Username)
	fmt.Print("--> " + user.Name)
	fmt.Print("--> " + user.Uid)
	cmd := exec.Command("bash", "-c", "sudo chown -R "+user.Username+":pkcs11 "+path+"/softoken")

	err22 := cmd.Run()
	if err22 != nil {

		fmt.Println(err22)
	}
}

func InstallSoftToken() error {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	cmd := exec.Command("bash", "-c", "cd src/ && sudo ./bootstrap.sh && sudo ./configure --prefix="+path+"/softoken && sudo make && sudo make install && sudo ldconfig && cd ../ && sudo chown "+user.Username+":pkcs11 -R softoken/")

	// Crear pipe para capturar commando en vivo

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t> %s\n", scanner.Text())
			ExecutionLog.Printf("\t%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return err
}

func StartToken() {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	_, err2 := exec.Command("bash", "-c", "pidof pkcsslotd").Output()
	if err2 != nil {
		cmd, err3 := exec.Command("bash", "-c", "sudo "+path+"/softoken/sbin/pkcsslotd").Output()
		if err3 != nil {
			fmt.Print("--> ERROR INICIANDO VIRTUAL ETOKEN ", cmd)
			os.Exit(1)
		} else {
			fmt.Print("--> VIRTUAL ETOKEN INICIADO CORRECTAMENTE.\n")
		}

	} else {
		fmt.Print("--> ERROR VIRTUAL TOKEN YA ESTA INCIADO EJECUTE -stop PARA DETENERLO.\n")
		os.Exit(1)
	}

}

func StopToken() {

	_, err2 := exec.Command("bash", "-c", "pidof pkcsslotd").Output()
	if err2 == nil {
		cmd, err3 := exec.Command("bash", "-c", "sudo pkill pkcsslotd").Output()
		if err3 != nil {
			fmt.Print("--> ERROR DETENIENDO VIRTUAL ETOKEN ", cmd)
			os.Exit(1)
		} else {
			fmt.Print("--> VIRTUAL ETOKEN DETENIDO CORRECTAMENTE.\n")
		}

	} else {
		fmt.Print("--> ERROR VIRTUAL TOKEN NO ESTA INCIADO EJECUTE -start PARA INICIARLO.\n")
		os.Exit(1)
	}

}

func InitToken() {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	_, oerr := exec.Command("bash", "-c", "pidof pkcsslotd").Output()
	if oerr != nil {

		fmt.Printf("--> DEBE INICIAR VIRTUAL ETOKEN ANTES DE INIT\n")
		os.Exit(1)

	}

	if _, err := os.Stat(path + "/softoken/init"); os.IsNotExist(err) {
		_ = os.Mkdir(path+"/softoken/init", os.ModePerm)

		_, err22 := exec.Command("bash", "-c", "printf 'virtualtoken - mfvilche' | "+path+"/softoken/sbin/pkcsconf -I -c 3 -S 87654321 -u -n 123456").Output()
		if err22 != nil {
			log.Fatal(err22)
		} else {
			fmt.Printf("--> TOKEN INIT OK\n")
		}

		_, err2 := exec.Command("bash", "-c", "openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout "+path+"/softoken/init/private.key -out "+path+"/softoken/init/cert.pem -subj '/CN=testing'").Output()
		if err2 != nil {
			log.Fatal(err2)
		} else {
			fmt.Printf("--> CERTIFICADOS AUTOFIRMADO GENERADO EN TOKEN CORRECTAMENTE\n")
		}

		_, err3 := exec.Command("bash", "-c", "openssl x509 -inform pem -outform der -in "+path+"/softoken/init/cert.pem -out "+path+"/softoken/init/cert.der").Output()
		if err3 != nil {
			log.Fatal(err3)
		} else {
			fmt.Printf("--> DER CERTIFICADO OK\n")
		}

		_, err4 := exec.Command("bash", "-c", "openssl rsa -inform pem -outform der -in "+path+"/softoken/init/private.key -out "+path+"/softoken/init/key.der").Output()
		if err4 != nil {
			log.Fatal(err4)
		} else {
			fmt.Printf("--> DER CLAVE PRIVADA OK\n")
		}

		_, err5 := exec.Command("bash", "-c", "pkcs11-tool --module "+path+"/softoken/lib/opencryptoki/libopencryptoki.so --slot 3 --login --pin 123456 --write-object "+path+"/softoken/init/key.der --type privkey --id 12 --label 'Testing'").Output()
		if err5 != nil {
			log.Fatal(err5)
		} else {

			fmt.Printf("--> CLAVE PRIVADA IMPORTADA EN TOKEN CORRECTAMENTE\n")
		}

		_, err6 := exec.Command("bash", "-c", "pkcs11-tool --module "+path+"/softoken/lib/opencryptoki/libopencryptoki.so --slot 3 --login --pin 123456 --write-object "+path+"/softoken/init/cert.der --type cert --id 12 --label 'Testing'").Output()
		if err6 != nil {
			log.Fatal(err6)
		} else {
			fmt.Printf("--> CERTIFICADOS IMPORTADO EN TOKEN CORRECTAMENTE\n")
		}

		out6, err6 := exec.Command("bash", "-c", "pkcs11-tool --module "+path+"/softoken/lib/opencryptoki/libopencryptoki.so --slot 3 --list-objects --type cert").Output()
		if err6 != nil {
			log.Fatal(err2)
		} else {
			fmt.Print(string(out6))
			fmt.Printf("--> CERTIFICADOS GENERADOS CORRECTAMENTE\n")
			fmt.Printf("--> ACCEDIENDO AL ETOKEN PARA VERIFICACION\n")
			_, err66 := exec.Command("bash", "-c", "pkcs11-tool --module "+path+"/softoken/lib/opencryptoki/libopencryptoki.so -L --slot 3 --login --pin 123456 -t -v").Output()
			if err6 != nil {
				log.Fatal(err66)
			} else {
				fmt.Print("--> ACCESO CORRECTO\n")
				fmt.Print("------------------------------------------->\n")
				fmt.Print("--> EL PIN DE SU VIRTUAL ETOKEN ES 123456\n")
				fmt.Print("--> EL SLOT DE SU VIRTUAL ETOKEN ES 3\n")
				fmt.Print("--> EL DRIVER DE SU ETOKEN ES " + path + "/softoken/lib/opencryptoki/libopencryptoki.so \n")
				fmt.Print("--> UTILICE pkcs11-tools PARA ACCEDER AL TOKEN.  pkcs11-tool --module " + path + "/softoken/lib/opencryptoki/libopencryptoki.so -L --slot 3 --login --pin 123456 \n")
				fmt.Print("------------------------------------------->\n")
			}

		}

	} else {
		fmt.Print("--> VIRTUAL TOKEN YA INICIALIZADO.\n")
		os.Exit(1)
	}

}

func AskSudo() {

	cmd22 := exec.Command("sudo", "id")
	cmd22.Stderr = os.Stderr
	cmd22.Stdin = os.Stdin
	cmd22.Stdout = os.Stdout
	err22 := cmd22.Run()
	if err22 != nil {
		fmt.Println(err22)
	}
}

func checkRoot() {

	cmd := exec.Command("id", "-u")
	output, _ := cmd.Output()
	i, _ := strconv.Atoi(string(output[:len(output)-1]))

	if i == 0 {
		fmt.Print("--> ERROR - NO ESTA PERMITIDO INICIAR EL PROGRAMA CON EL USUARIO ROOT.\n")
		os.Exit(1)
	}
}
