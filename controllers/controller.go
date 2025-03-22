package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/gofiber/fiber/v2"
)

func FileUpload(c *fiber.Ctx)error {
	f, err := c.FormFile("file")
	if err!=nil{
		return err
	}
	fmt.Println(f.Filename)
	file, err := f.Open() 
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf,file); err!=nil{
		return err
	}
	expiryTime := time.Now().UTC().Add(1 * time.Hour)
	UploadFileToStorage(buf.Bytes(),f.Filename)
	url := getUrl(f.Filename,expiryTime)
	return c.Status(http.StatusOK).JSON(fiber.Map{"fileUrl":url})

}


func UploadFileToStorage(buf []byte, filename string){
	client := getAzureClient()
	blobData := string(buf)
	blobContainerName := "files"
	expiryTime := time.Now().UTC().Add(1 * time.Hour)
	uploadfileResp, err := client.UploadStream(context.TODO(),
		blobContainerName,
		filename,
		strings.NewReader(blobData),
		&azblob.UploadStreamOptions{
			Tags: map[string]string{"ExpiryTime":expiryTime.String() },
		})
	fmt.Println(uploadfileResp)
	handleError(err)


}



func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getAzureClient() *azblob.Client {
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		log.Fatal("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}
	accountName := "securesharesta"

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	handleError(err)


	return client
}

func getUrl(filename string,expTime time.Time) string{
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		log.Fatal("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}
	accountName := "securesharesta"

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	client , err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	handleError(err)
	_ = client // just for holding the variable client without using it
	sasQueryParameters,err := sas.BlobSignatureValues{
		ExpiryTime: expTime,
		ContainerName: "files",
		Permissions:  to.Ptr(sas.ContainerPermissions{Read: true}).String(),
		BlobName: filename,
	}.SignWithSharedKey(cred)
	sasUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",accountName,"files",filename,sasQueryParameters.Encode())

	return sasUrl

}	