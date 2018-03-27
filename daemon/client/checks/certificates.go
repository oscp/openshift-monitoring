package checks

import (
	"encoding/pem"
	"log"
	"os"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"crypto/x509"
	"time"
	"gopkg.in/yaml.v2"
	"encoding/base64"
	"errors"
	"net/http"
	"crypto/tls"
)

type Cert struct {
	File string
	DaysLeft int
}

type KubeConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster   string `yaml:"cluster"`
			Namespace string `yaml:"namespace"`
			User      string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func decodeCertBlocks(data []byte) []*pem.Block {
	var blocks []*pem.Block
	block, rest := pem.Decode([]byte(data))

	if block != nil {
		blocks = append(blocks, block)
	}

	if len(rest) > 0 {
		return append(blocks, decodeCertBlocks(rest)...)
	} else {
		return blocks
	}
}

func getCertFiles(filePaths []string) (error, []string) {
	var certFiles []string

	for _, path := range filePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("Path %s does not exist.", path)
			continue
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			msg := fmt.Sprintf("could not read directory %s (%s)", path, err.Error())
			log.Println(msg)
			return errors.New(msg), nil
		}

		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".crt" {
				continue
			}

			certFiles = append(certFiles, filepath.Join(path, file.Name()))
		}
	}

	return nil, certFiles
}

func getExpiredCerts(filePaths []string, days int) (error, []Cert) {
	var expiredCerts []Cert

	err, certFiles := getCertFiles(filePaths)
	if err != nil {
		return errors.New(fmt.Sprintf("could not get files (%s)", err.Error())), nil
	}

	for _, file := range certFiles {

		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.New(fmt.Sprintf("could not read file %s", file)), nil
		}

		blocks := decodeCertBlocks([]byte(data))

		for _, block := range blocks {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return errors.New(fmt.Sprintf("certificate parsing error (%s)", err.Error())), nil
			}

			daysLeft := cert.NotAfter.Sub(time.Now()).Hours() / 24

			if int(daysLeft) <= days {
				log.Println(fmt.Sprintf("%s expires in %d days", file, int(daysLeft)))
				expiredCerts = append(expiredCerts, Cert { File: file, DaysLeft: int(daysLeft) })
			}
		}

	}

	return nil, expiredCerts
}

func CheckUrlSslCertificates(urls []string, days int) error {
	log.Printf("Checking expiry date for SSL certificates (%d days) via urls.", days)

	var certErrorList []string

	for _, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			msg := fmt.Sprintf("creating request failed for %s (%s)", url, err.Error())
			log.Println(msg)
			return errors.New(msg)
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		hc := &http.Client{Transport: tr}

		resp, err := hc.Do(req)
		if err != nil {
			msg := fmt.Sprintf("get request failed for %s (%s)", url, err.Error())
			log.Println(msg)
			return errors.New(msg)
		}

		if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
			for _, cert := range resp.TLS.PeerCertificates {
				daysLeft := cert.NotAfter.Sub(time.Now()).Hours() / 24

				if int(daysLeft) <= days {
					msg := fmt.Sprintf("certificate %s from %s expires in %d days", cert.Subject, url, int(daysLeft))
					log.Println(msg)
					certErrorList = append(certErrorList, msg)
				}
			}
		}
	}

	if len(certErrorList) > 0 {
		var errorMessage string
		for _, msg := range certErrorList {
			errorMessage = errorMessage + msg + " "
		}
		return errors.New(errorMessage)
	}

	return nil
}

func CheckFileSslCertificates(filePaths []string, days int) error {
	log.Printf("Checking expiry date for SSL certificates (%d days) in files.", days)

	var certErrorList []string

	err, expiredCerts := getExpiredCerts(filePaths, days)
	if err != nil {
		msg := "could not get expired certificates"
		log.Println(msg)
		return errors.New(msg)
	}

	for _, expiredCert := range expiredCerts {
		certErrorList = append(certErrorList, fmt.Sprintf("%s expires in %d days", expiredCert.File, expiredCert.DaysLeft))
	}

	if _, err := os.Stat("/root/.kube/config"); os.IsNotExist(err) {
		log.Println("File /root/.kube/config does not exist.")
	} else {

		data, err := ioutil.ReadFile("/root/.kube/config")
		if err != nil {
			msg := "could not read file /root/.kube/config"
			log.Println(msg)
			return errors.New(msg)
		}

		var kubeConfig KubeConfig

		err = yaml.Unmarshal(data, &kubeConfig)
		if err != nil {
			msg := fmt.Sprintf("unmarshalling /root/.kube/config failed (%s)", err.Error())
			log.Println(msg)
			return errors.New(msg)
		}

		for _, cluster := range kubeConfig.Clusters {
			if len(cluster.Cluster.CertificateAuthorityData) > 0 {

				certBytes, err := base64.StdEncoding.DecodeString(cluster.Cluster.CertificateAuthorityData)
				if err != nil {
					msg := fmt.Sprintf("can't base64 decode cert (%s)", err.Error())
					log.Println(msg)
					return errors.New(msg)
				}

				block, _ := pem.Decode(certBytes)

				cert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					msg := fmt.Sprintf("certificate parsing error (%s)", err.Error())
					log.Println(msg)
					return errors.New(msg)
				}

				daysLeft := cert.NotAfter.Sub(time.Now()).Hours() / 24

				if int(daysLeft) <= days {
					msg := fmt.Sprintf("certificate-authority-data from /root/.kube/config expires in %d days", int(daysLeft))
					log.Println(msg)
					certErrorList = append(certErrorList, msg)
				}
			}
		}

		for _, user := range kubeConfig.Users {
			if len(user.User.ClientCertificateData) > 0 {

				certBytes, _ := base64.StdEncoding.DecodeString(user.User.ClientCertificateData)
				block, _ := pem.Decode(certBytes)

				cert, err := x509.ParseCertificate(block.Bytes)

				if err != nil {
					msg := fmt.Sprintf("certificate parsing error (%s)", err.Error())
					log.Println(msg)
					return errors.New(msg)
				}

				daysLeft := cert.NotAfter.Sub(time.Now()).Hours() / 24

				if int(daysLeft) <= days {
					msg := fmt.Sprintf("client-certificate-data from /root/.kube/config expires in %d days", int(daysLeft))
					log.Println(msg)
					certErrorList = append(certErrorList, msg)
				}
			}
		}

	}

	if len(certErrorList) > 0 {
		var errorMessage string
		for _, msg := range certErrorList {
			errorMessage = errorMessage + msg + " "
		}
		return errors.New(errorMessage)
	}

	return nil
}
