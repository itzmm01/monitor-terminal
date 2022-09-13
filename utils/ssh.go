package utils

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

// ClientConfig 连接的配置
type ClientConfig struct {
	Host       string //ip
	Port       int64  // 端口
	Username   string //用户名
	Password   string //密码
	KeyFile    string //密钥文件
	Timeout    time.Duration
	sshClient  *ssh.Client  //ssh client
	sftpClient *sftp.Client //sftp client
	LastResult string       //最近一次运行的结果
}

func publicKeyAuthFunc(keyPath string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		Log.Println("Failed to read ssh key file:", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		Log.Println("Failed to signature ssh key file: ", err)
	}
	return ssh.PublicKeys(signer)
}

// CreateClient 建立连接
func (cliConf *ClientConfig) CreateClient() error {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		err        error
	)
	if cliConf.Timeout == 0 {
		cliConf.Timeout = 5
	}
	config := ssh.ClientConfig{
		User: cliConf.Username,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: cliConf.Timeout * time.Second,
	}
	if cliConf.KeyFile == "" {
		config.Auth = []ssh.AuthMethod{ssh.Password(cliConf.Password)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(cliConf.KeyFile)}
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	if sshClient, err = ssh.Dial("tcp", addr, &config); err != nil {
		return err
	}
	cliConf.sshClient = sshClient

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return err
	}
	cliConf.sftpClient = sftpClient
	return nil
}

// RunShell 执行命令
func (cliConf *ClientConfig) RunShell(shell string) (res string, error1 error) {
	var (
		session *ssh.Session
		err     error
	)
	//获取session，这个session是用来远程执行操作的
	if session, err = cliConf.sshClient.NewSession(); err != nil {
		return "", err
	}
	defer session.Close()
	//执行shell
	if output, err := session.CombinedOutput(shell); err != nil {
		return "", err
	} else {
		if len(output) > 0 {
			lastStr := output[len(output)-1]
			if lastStr == byte(10) {
				cliConf.LastResult = string(output[0 : len(output)-1])
			} else {
				cliConf.LastResult = string(output)
			}
		} else {
			cliConf.LastResult = string(output)
		}
	}
	result := cliConf.LastResult
	return result, nil
}

// Upload 上传文件
func (cliConf *ClientConfig) Upload(srcPath, dstPath string) error {
	srcFile, _ := os.Open(srcPath) //本地
	destDir := path.Dir(dstPath)
	cliConf.RunShell("mkdir -p " + destDir)
	dstFile, _ := cliConf.sftpClient.Create(dstPath) //远程
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}
		_, _ = dstFile.Write(buf[:n])
	}
	return nil
}

// UploadDirectory 上传目录
func (cliConf *ClientConfig) UploadDirectory(srcDir, dstPath string) error {
	srcFiles, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, backupDir := range srcFiles {
		srcFilePath := path.Join(srcDir, backupDir.Name())
		dstFilePath := path.Join(dstPath, backupDir.Name())
		if backupDir.IsDir() {
			cliConf.sftpClient.Mkdir(dstFilePath)
			cliConf.UploadDirectory(srcFilePath, dstFilePath)
		} else {
			cliConf.Upload(srcFilePath, dstFilePath)
		}
	}
	return nil
}

// Download 下载文件
func (cliConf *ClientConfig) Download(srcPath, dstPath string) error {
	srcFile, _ := cliConf.sftpClient.Open(srcPath) //远程
	dstFile, _ := os.Create(dstPath)               //本地
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	return nil
}

// DownloadDirectory 下载目录
func (cliConf *ClientConfig) DownloadDirectory(srcPath, dstPath string) error {
	w := cliConf.sftpClient.Walk(srcPath)
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		fileName := strings.Split(w.Path(), srcPath)
		stat, _ := cliConf.sftpClient.Stat(w.Path())
		if stat.IsDir() {
			err := os.MkdirAll(dstPath+fileName[len(fileName)-1], 0755)
			if err != nil {
				return err
			}
		} else {
			err := cliConf.Download(w.Path(), dstPath+fileName[len(fileName)-1])
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// Delete 删除远程目录
func (cliConf *ClientConfig) Delete(filePath string) error {
	return cliConf.sftpClient.Remove(filePath)
}
