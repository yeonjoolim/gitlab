package main

import (
				"fmt"
				"archive/tar"
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
		defer os.RemoveAll(tmpPath)
		if result == false  {
			os.Exit(1)
		}
		str := read_mapping_signed(os.Args[1]+"-sign.gob")
		result = server_verify("./dev1/dev1.crt",str)	

		if result == false  {
			os.Exit(1)
		} 
}

func server_verify(keypath string, strl map[string]string) bool{
		var Pubkey *rsa.PublicKey

		Pubkey = Getpubcrt(keypath)
		
		veri := true
        ctr := 0
		i := 0
	
		for key, val := range strl {
			for i=0; i<len(layerIds); i++{
				if key == layerIds[i]{
					result := layer_verify(Pubkey, val, GetHash(layersPath[i]))
					if result == false {
						fmt.Printf("Verify Fail %s",layerIds[i])
						veri = false
					}else{
						ctr++
					}
				}	
			}
			if ctr == len(layerIds){
				veri = true
			}
		}
		if veri == true{
			fmt.Printf("verify success\n")
			return true
		} else {
			fmt.Printf("Verify fail\n")
			return false
		}
}

func GetHash(path string) []byte{
		file, err := os.Open(path)
        if err != nil {
                log.Fatal("ExtractTarGz: NewReader failed")
        }

        defer file.Close()

        sha := sha256.New()
        if _, err := io.Copy(sha, file); err != nil {
                log.Fatal(err)
        }
        h2 := sha.Sum(nil)

		return h2
}

func layer_verify(Pubkey *rsa.PublicKey, val string, hash []byte) bool{
			decoded_sig, err := hex.DecodeString(val)
				
			veri := Verify(Pubkey,hash,decoded_sig)
				
			if err != nil {
				log.Fatal(err)
			}
			//log.Printf("sign data: %s\n",hash)
			//log.Printf("signed data: %s\n",decoded_sig)
			return veri
}

func read_mapping_signed(path string) map[string]string{
    // Open a RO file
    decodeFile, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer decodeFile.Close()

    // Create a decoder
    decoder := gob.NewDecoder(decodeFile)

    // Place to decode into
    il := make(map[string]string)

    // Decode -- We need to pass a pointer otherwise accounts2 isn't modified
    decoder.Decode(&il)
    
	Signed_data := make(map[string]string)

	for key,val := range il {
        Signed_data[key] = strings.TrimPrefix(val,"%!(EXTRA string")
    }
	return Signed_data
}

func image_layer_load(imageName string) bool {
		//Create a temporary folder where the docker image layers are going to be stored
		tmpPath := createTmpPath(tmpPrefix)
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

func Getpubcrt(path string) *rsa.PublicKey{
        fd, _ := os.Open(path)
        defer fd.Close()
        stat, _ := fd.Stat()
        var buf = make([]byte, stat.Size())
        _, _ = fd.Read(buf)
        data, _  := pem.Decode(buf)
        var cert* x509.Certificate
        cert, _ = x509.ParseCertificate(data.Bytes)
        pubkey := cert.PublicKey.(*rsa.PublicKey)
        return pubkey
}


func Verify(pub *rsa.PublicKey, hash []byte, sig []byte) bool{
        err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash, sig)
    if (err != nil) {
                return false
    } else {
                return true
    }
}
