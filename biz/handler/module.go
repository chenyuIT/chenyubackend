package handler

import (
	"chenyucms/biz/util"
	"chenyucms/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/goravel/framework/facades"
	"google.golang.org/grpc/codes"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Client call
type respon struct {
	code codes.Code
	data any
}

// Module .
func Module(ctx context.Context, c *app.RequestContext) {

	//Method of Module
	method := string(c.Method())
	//Module name of mod
	param := c.Param("mod")
	//Module Service of serv
	serv := c.Param("serv")
	//Header of Authorization
	header := c.GetHeader("Authorization")
	//Request Data
	dat, _ := c.MultipartForm()
	data := dat.Value
	//for display format
	var mydata, _ = json.Marshal(data)

	//install  module
	if param == "install" {

		if ModuleInstall(serv) {
			c.JSON(consts.StatusOK, utils.H{
				"message": "install module complete",
			})
		} else {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"message": "install module Failed",
			})
		}
		return
	}

	//update module
	if param == "update" {
		if ModuleUpdate(serv) {
			c.JSON(consts.StatusOK, utils.H{
				"message": "update module complete",
			})
		} else {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"message": "update module Failed",
			})
		}
		return
	}

	//delete module
	if param == "delete" {
		if ModuleDelete(serv) {
			c.JSON(consts.StatusOK, utils.H{
				"message": "delete module complete",
			})
		} else {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"message": "delete module Failed",
			})
		}
		return
	}

	//start module
	if param == "start" {
		if ModuleStart(serv) {
			c.JSON(consts.StatusOK, utils.H{
				"message": "Start module complete",
			})
		} else {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"message": "Start module Failed",
			})
		}
		return
	}

	//stop module
	if param == "stop" {
		if ModuleStop(serv) {
			c.JSON(consts.StatusOK, utils.H{
				"message": "Stop module complete",
			})
		} else {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"message": "Stop module Failed",
			})
		}
		return
	}
	//make a client call
	modlist := facades.Cache().Store("memory").Get(param, config.Module{})
	//断言
	if mod, ok := modlist.(config.Module); ok {
		fmt.Printf("Name:%s,  Host:%s,  Port:%d", mod.Name, mod.IP, mod.Port)
		resp := ClientRequest(mod.IP, int(mod.Port), serv, method, mydata)
		//output the result for display
		fmt.Println("返回内容:", resp.data)
		c.JSON(consts.StatusOK, resp.data)

	} else {
		//output the result for display
		c.JSON(consts.StatusBadRequest, utils.H{
			"message": "Module-test:" + method + "  Param:" + param + "  Service:" + serv + "  Head_authorization:" + string(header) + "  Req:" + string(mydata),
		})

	}
}

// ModuleInstall install module process
func ModuleInstall(srv string) bool {
	//check module zip package existed
	file := srv + ".zip"
	if facades.Storage().Disk("uploads").Exists(file) {
		//unzip srv.zip
		util.Unzip(facades.Storage().Disk("uploads").Path("./")+file, "module/"+srv)
		if facades.Storage().Disk("module").Exists("/"+srv+"/"+srv) || facades.Storage().Disk("module").Exists("/"+srv+"/"+srv+".exe") {
			//unzip succeed
			if facades.Storage().Disk("module").Exists("/" + srv + "/" + srv) {
				fmt.Println("Linux p")
			} else {
				fmt.Println("Windows p")
			}
			return true
		} else {
			//unzip failed
			fmt.Println("File unzip failed!  ", file)
		}

	} else {
		fmt.Println("File not Exist!  ", file)
	}
	//unzip module zip
	//check .env is valid or invalid
	//check the module path exist or not exist
	//--if not exist then make fold and copy files
	//--if exist then show notice for update not install
	return false
}

// ModuleUpdate update module process
func ModuleUpdate(srv string) bool {
	return false
}

// ModuleDelete delete module process
func ModuleDelete(srv string) bool {
	err := facades.Storage().DeleteDirectory("../../module/" + srv)
	if err == nil {
		return true
	}
	return false
}

// ModuleStart start module process
func ModuleStart(srv string) bool {
	operatingSystem := runtime.GOOS
	wd, _ := os.Getwd()
	execPath := wd + "\\module\\" + srv + "\\" // exec command path
	curhost := "0.0.0.0"
	curport := facades.Config().GetInt("module.module_port_current", 0)

	port, err := CheckPortAvailable(curport)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	} else {
		fmt.Printf("Available port: %d\n", port)
		curport = port
	}

	cmd := exec.Command(execPath + srv) // set exec module command
	if setEnv(curhost, curport, execPath+".env") {
		if operatingSystem == "windows" {
			cmd = exec.Command("cmd", "/C", "start "+execPath+srv) // set exec module command
		} else if operatingSystem == "linux" {
			cmd = exec.Command("nohup", "", execPath+srv) // set exec module command
		} else {
			cmd = exec.Command(execPath + srv) // set exec module command
		}

		cmd.Dir = execPath //set exec Path
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(srv+" 执行命令时出错:", err)
			return false
		}
		fmt.Println(execPath + " 执行完成")
		//写入config.module
		facades.Config().Add("module.module_port_current", curport+1)
		facades.Config().Add("module.module_host_current", curhost)
		var mod config.Module
		mod = config.Module{Name: srv, IP: curhost, Port: int32(curport)}
		facades.Cache().Store("memory").Put(srv, mod, 0)
		modlist := facades.Cache().Store("memory").Get(srv, config.Module{})
		//断言
		if mymod, ok := modlist.(config.Module); ok {
			mod = mymod
		}
		//输出
		//fmt.Println(reflect.TypeOf(modlist))
		//fmt.Printf("Name:%s,  Host:%s,  Port:%d", mod.Name, mod.IP, mod.Port)

		return true
	} else {
		return false
	}

}

// ModuleStop Stop module process
func ModuleStop(srv string) bool {
	isRunning, err := isProcessRunning(srv)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}
	if isRunning {
		fmt.Printf("%s 程序正在运行\n", srv)
		err := stopProcess(srv)
		if err != nil {
			fmt.Printf("Error stopping process: %v\n", err)
			return false
		}
		fmt.Printf("%s 进程已停止\n", srv)
		return true
	} else {
		fmt.Printf("%s 程序未运行\n", srv)
	}
	return false
}

func isProcessRunning(processName string) (bool, error) {
	operatingSystem := runtime.GOOS
	cmd := exec.Command("ps", "aux") //default on linux
	if operatingSystem == "windows" {
		cmd = exec.Command("tasklist")
	} else if operatingSystem == "linux" {
		cmd = exec.Command("ps", "aux")
	} else {
		fmt.Println("无法识别的操作系统")
		return false, nil
	}

	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, processName) {
			return true, nil
		}
	}

	return false, nil
}

func stopProcess(processName string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// 在Windows下使用taskkill命令
		cmd = exec.Command("taskkill", "/F", "/IM", processName+".exe")
	case "linux":
		// 在Linux下使用pkill命令
		cmd = exec.Command("pkill", processName)
	default:
		return fmt.Errorf("Unsupported operating system: %s", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// 重置.env(Module的)有关APP_HOST和 APP_PORT的内容并将Module的设置写入 config.module.modlist配置中
func setEnv(host string, port int, file string) bool {
	envcontent, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("读取.env文件时出错：%v\n", err)
		return false
	}
	// 替换APP_HOST和APP_PORT的值
	newContent := replaceEnvVariables(string(envcontent), "APP_HOST", host)
	newContent = replaceEnvVariables(newContent, "APP_PORT", strconv.Itoa(port))
	// 将新内容写回.env文件
	err = ioutil.WriteFile(file, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("写入.env文件时出错：%v\n", err)
		return false
	}
	fmt.Println("已成功读取、替换并写入.env文件")
	return true
}

// 使用正则表达式替换.env文件中的环境变量值
func replaceEnvVariables(content, variableName, newValue string) string {
	re := regexp.MustCompile(fmt.Sprintf(`\b%s\b=.*`, variableName))
	return re.ReplaceAllString(content, fmt.Sprintf("%s=%s", variableName, newValue))
}

//检查端口是否可用

func CheckPortAvailable(startPort int) (int, error) {
	for port := startPort; port <= 65535; port++ {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found")
}

func ClientRequest(host string, port int, api string, method string, data []byte) respon {
	resp := respon{}
	c, err := client.NewClient()
	if err != nil {
		resp.code = codes.Unknown
		resp.data = err.Error()
		return resp
	}
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.SetMethod(method)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetBodyRaw(data)
	req.SetRequestURI("http://" + host + ":" + strconv.Itoa(port) + "/" + api)
	err = c.Do(context.Background(), req, res)
	if err != nil {
		resp.code = codes.Unknown
		resp.data = err.Error()
		return resp
	}
	fmt.Printf("调用返回:%v", string(res.Body()))
	resp.code = codes.OK
	resp.data = res.Body()
	return resp
}
