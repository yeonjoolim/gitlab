package main

import (
		"archive/tar"
		"fmt"
		"io"
		"io/ioutil"
		"path/filepath"
		"context"
		"encoding/json"
		"strings"
		"github.com/docker/docker/client"
		"os"
                "encoding/gob"
                "encoding/hex"
                "log"
                "crypto"
                "crypto/rand"
                "crypto/rsa"
                "crypto/sha256"
                "crypto/x509"
                "encoding/pem"
		)

const tmpPrefix="pel-lab-"

var layerIds []string
var layersPath []string
var tmpPath string

func main() {
		result := image_layer_load(os.Args[1])
		if result == false  {
			os.Exit(1)
		}
		defer os.RemoveAll(tmpPath)
		
		var encodesign []string
		for i := 0; i<len(layerIds); i++{
			encodesign = append(encodesign, layer_sign(layersPath[i],"./repo/repo.key"))
        }
        
	    mapping_sign(encodesign)
		fmt.Printf("Resign Success")
}

func filew(path string, data []byte) {
     fd, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
     defer fd.Close()
     _,_ = fd.Write([]byte(data))
}

func Getpri(path string) *rsa.PrivateKey {
        fd, _ := os.Open(path)
        defer fd.Close()
        stat, _ := fd.Stat()
        var buf = make([]byte, stat.Size())
        _, _ = fd.Read(buf)
        data, _ :=pem.Decode(buf)
        prikey, _ := x509.ParsePKCS1PrivateKey(data.Bytes)
        return prikey
}

func sign(pri *rsa.PrivateKey, hash []byte) []byte{
     sig, _ := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, hash)
     return sig
}

func layer_sign(path string, keypath string) string{
        var Server_sign []byte
        var Prikey *rsa.PrivateKey

        file, err := os.Open(path)
        if err != nil {
                log.Fatal("ExtractTarGz: NewReader failed")
        }

        Prikey = Getpri(keypath)
        defer file.Close()

        sha := sha256.New()
        if _, err := io.Copy(sha, file); err != nil {
                log.Fatal(err)
        }
        h2 := sha.Sum(nil)

        Server_sign = sign(Prikey,h2)
        encodesign := hex.EncodeToString(Server_sign)
        return encodesign

}

func mapping_sign(server_sign []string) {
	Signdata := make(map[string]string)
    	for i := 0; i<len(layerIds); i++{ 
		Signdata[layerIds[i]] = server_sign[i]
    	}
    	path := os.Args[1] + "-resign.gob"
    	create_signfile(path,Signdata)
    
}

func create_signfile(path string, Signdata map[string]string){
    encodeFile, err := os.Create(path)
    if err != nil {
        panic(err)
    }

    // Since this is a binary format large parts of it will be unreadable
    encoder := gob.NewEncoder(encodeFile)

    // Write to the file
    if err := encoder.Encode(Signdata); err != nil {
        panic(err)
    }
    encodeFile.Close()
    log.Printf("Resigned signature data is in %s\n",path)
}

func image_layer_load(imageName string) bool {
		//Create a temporary folder where the docker image layers are going to be stored
		tmpPath = createTmpPath(tmpPrefix)
		
		saveDockerImage(imageName, tmpPath)
		getImageLayerIds(tmpPath)
		
		//print_layerinfo()
	
		return true
}

func print_layerinfo() {
        for i := 0; i < len(layerIds); i++ {
                fmt.Printf("Layer %d %s",i, layerIds[i])
                fmt.Printf(" in %s\n",layersPath[i])
        }
}

// TODO Add support for older version of docker

type manifestJSON struct {
	Layers []string
}

// saveDockerImage saves Docker image to temorary folder
func saveDockerImage(imageName string, tmpPath string) {
		docker := createDockerClient()

		imageReader, err := docker.ImageSave(context.Background(), []string{imageName})
		if err != nil {
			fmt.Printf("Could not save Docker image [%s]: %v", imageName, err)
		}
		defer imageReader.Close()
		if err = untar(imageReader, tmpPath); err != nil {
			fmt.Printf("Could not save Docker image: could not untar [%s]: %v", imageName, err)
		}
}

func createDockerClient() client.APIClient {
	docker, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		fmt.Printf("Could not create a Docker client: %v", err)
	}
	return docker
}

// getImageLayerIds reads LayerIDs from the manifest.json file
func getImageLayerIds(path string) {
	manifest := readManifestFile(path)

	for _, layer := range manifest[0].Layers {
		  layersPath = append(layersPath, path+"/"+layer)
		  layerIds = append(layerIds, strings.TrimSuffix(layer, "/layer.tar"))
	}
}

// readManifestFile reads the local manifest.json
func readManifestFile(path string) []manifestJSON {
	manifestFile := path + "/manifest.json"
    mf, err := os.Open(manifestFile)
	if err != nil {
	  fmt.Printf("Could not read Docker image layers: could not open [%s]: %v", manifestFile, err)
	}
	//fmt.Printf("%#v\n",mf)
	defer mf.Close()

    return parseAndValidateManifestFile(mf)
}

// parseAndValidateManifestFile parses the manifest.json file and validates it
func parseAndValidateManifestFile(manifestFile io.Reader) []manifestJSON {
	var manifest []manifestJSON
		if err := json.NewDecoder(manifestFile).Decode(&manifest); err != nil {
			fmt.Printf("Could not read Docker image layers: manifest.json is not json: %v", err)
		} else if len(manifest) != 1 {
			fmt.Printf("Could not read Docker image layers: manifest.json is not valid")
		} else if len(manifest[0].Layers) == 0 {
			fmt.Printf("Could not read Docker image layers: no layers can be found")
		}
	//fmt.Printf("%#v\n",manifest)
	return manifest
}

func createTmpPath(tmpPrefix string) string {
	tmpPath, err := ioutil.TempDir("", tmpPrefix)
		if err != nil {
			fmt.Printf("Could not create temporary folder: %s", err)
		}
	return tmpPath
}

// untar uses a Reader that represents a tar to untar it on the fly to a target folder
func untar(imageReader io.ReadCloser, target string) error {
	//fmt.Println(target)
	tarReader := tar.NewReader(imageReader)

			   for {
				   header, err := tarReader.Next()
					   if err == io.EOF {
						   break
					   } else if err != nil {
						   return err
					   }

	path := filepath.Join(target, header.Name)
		  if !strings.HasPrefix(path, filepath.Clean(target) + string(os.PathSeparator)) {
			  return fmt.Errorf("%s: illegal file path", header.Name)
		  }
	info := header.FileInfo()
		  if info.IsDir() {
			  if err = os.MkdirAll(path, info.Mode()); err != nil {
				  return err
			  }
			  continue
		  }

	  file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		  if err != nil {
			  return err
		  }
	  defer file.Close()
		  if _, err = io.Copy(file, tarReader); err != nil {
			  return err
		  }
			   }
		   return nil
}

