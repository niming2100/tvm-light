package contract

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	tvm_conf "tvm-light/config"
)

type Contract struct {
	peerAddress     string
	contractName    string
	contractType    string
	contractPath    string
	contractVersion string
	channelID       string
	orgName         string
	args            string
	action          string
}

type Args struct {
	Args interface{} `json:"Args"`
}

// {\"Args\":[\"query\",\"a\"]}
const (
	CMD_DOCKER = "docker"
)

var docker_command []string = []string{"exec", "cli"}

func NewContract(peerAddress string, contractName string, contractType string, contractPath string, contractVersion string, channelID string, orgName string, args string, action string) *Contract {
	return &Contract{peerAddress: peerAddress, contractName: contractName, contractType: contractType, contractPath: contractPath, contractVersion: contractVersion, channelID: channelID, orgName: orgName, args: args, action: action}
}

func (c *Contract) InstallContract() ([]byte,error) {
	var filePath string = tvm_conf.GetDockerPath() + c.contractPath[len(tvm_conf.GetContractPath()):];
	params := []string{"peer", "chaincode", "install", "-n", c.contractName, "-p", filePath, "-v", c.contractVersion};
	cmd := exec.Command(CMD_DOCKER, append(docker_command, params...)...);
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	stderrOut, err := cmd.StderrPipe()
	if err != nil {
		return nil,err
	}
	// 保证关闭输出流
	defer stdout.Close()
	defer stderrOut.Close()
	// 运行命令
	fmt.Println(cmd.Args)
	if err := cmd.Start(); err != nil {
		return nil,err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	errBytes, err := ioutil.ReadAll(stderrOut)
	if err != nil {
		return nil,err
	}
	fmt.Println(string(errBytes))
	fmt.Println(string(opBytes))
	cmd.Wait()
	return opBytes,nil;
}

func (c *Contract) RunContract() ([]byte, error) {

	var resp []byte;
	switch c.action {
	case "instantiate":
		return c.instantiate()
		break;
	case "install":
		return c.InstallContract()
	default:
		return c.execute()
		break;
	}
	return resp, nil;
}

func (c *Contract) instantiate() ([]byte, error) {
	params := []string{"peer", "chaincode", "instantiate", "-o", tvm_conf.GetOrderServer(), "-C", c.channelID, "-n", c.contractName, "-v", c.contractVersion, "-c", c.args};
	cmd := exec.Command(CMD_DOCKER, append(docker_command, params...)...);
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	stderrOut, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 保证关闭输出流
	defer stdout.Close()
	defer stderrOut.Close()
	// 运行命令
	fmt.Println(cmd.Args)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	errBytes, err := ioutil.ReadAll(stderrOut)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	cmd.Wait()
	fmt.Println(string(errBytes))
	fmt.Println(string(opBytes))
	return opBytes, err;
}

func (c *Contract) execute() ([]byte, error) {
	// peer chaincode $action -C $channelID -n $cname -c $args
	params := []string{"peer", "chaincode", c.action, "-C", c.channelID, "-n", c.contractName, "-c", c.args};
	cmd := exec.Command(CMD_DOCKER, append(docker_command, params...)...);
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	stderrOut, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 保证关闭输出流
	defer stdout.Close()
	defer stderrOut.Close()
	// 运行命令
	fmt.Println(cmd.Args)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)

	errBytes, err := ioutil.ReadAll(stderrOut)

	fmt.Println(string(errBytes))
	fmt.Println(string(opBytes))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	cmd.Wait()
	return opBytes, err
}