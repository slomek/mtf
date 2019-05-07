package components

import (
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type MySQL struct {
	Pass     string
	Port     string
	DB       []string
	Hostname string
	Network  string
}

func (c *MySQL) Start() {
	time.Sleep(time.Second * 5)
	cmd := `docker run --rm -d --network=mtf_net --name mysql_mtf --hostname=mysql_mtf --env MYSQL_ROOT_PASSWORD=test --env MYSQL_DATABASE=test_db -p 3306:3306 mysql --default-authentication-plugin=mysql_native_password`
	run(cmd)
}

func (c *MySQL) Stop() {
	run("docker stop mysql_mtf")
}

func (c *MySQL) Ready() {
	waitForOpenPort("localhost", "3306")
}

func run(s string) error {
	args := strings.Split(s, " ")
	cmd := exec.Command(args[0], args[1:len(args)]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Printf("[ERROR] cmd start: %v , \n", err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("[ERROR] cmd wait: %v , \n", err)
		return err
	}
	//log.Printf("[DEBUG] output: %v\n", string(buff))
	return nil
}

func waitForOpenPort(host, port string) {
	firstRun := true
	for {
		if firstRun {
			firstRun = false
		} else {
			time.Sleep(time.Millisecond * 50)
		}
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Millisecond*500)
		if err != nil {
			continue
		}
		buff := make([]byte, 100)
		if _, err = conn.Read(buff); err != nil {
			conn.Close()
			continue
		}
		conn.Close()
		return
	}
}
