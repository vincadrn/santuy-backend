package oci

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"vincadrn.com/santuy/internal/model"
)

var (
	obsClient       objectstorage.ObjectStorageClient
	NAMESPACE       string = os.Getenv("OCI_NAMESPACE")
	BUCKETNAME      string = os.Getenv("OCI_BUCKET_NAME")
	OBJECTENVPREFIX string = os.Getenv("OCI_OBJECT_PREFIX")
)

func init() {
	tenancyString := os.Getenv("OCI_TENANCY")
	userString := os.Getenv("OCI_USER")
	region := os.Getenv("OCI_REGION")
	fingerprint := os.Getenv("OCI_FINGERPRINT")
	privateKey := strings.ReplaceAll(os.Getenv("OCI_PRIVATE_KEY"), "\\n", "\n")

	slog.Info("Authenticating to OCI", "tenancy", tenancyString)
	slog.Info("Authenticating to OCI", "user", userString)
	slog.Info("Authenticating to OCI", "region", region)
	slog.Info("Authenticating to OCI", "fingerprint", fingerprint)

	provider := common.NewRawConfigurationProvider(
		tenancyString,
		userString,
		region,
		fingerprint,
		privateKey,
		nil,
	)

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		slog.Error("Cannot create object storage client")
		return
	}

	obsClient = client
}

func GetClient() *objectstorage.ObjectStorageClient {
	return &obsClient
}

func generateObjectPrefix(domainPrefix string, id string) string {
	objectNameBuilder := strings.Builder{}
	objectNameBuilder.WriteString(OBJECTENVPREFIX)
	objectNameBuilder.WriteString("/")
	objectNameBuilder.WriteString(domainPrefix)
	objectNameBuilder.WriteString("=")
	objectNameBuilder.WriteString(id)
	objectNameBuilder.WriteString("/")

	return objectNameBuilder.String()
}

func generateRandomSuffixWithFormat(format string) string {
	randomSuffixBuilder := strings.Builder{}
	randomSuffixBuilder.WriteString(strconv.Itoa(rand.Int()))
	randomSuffixBuilder.WriteString(strconv.Itoa(rand.Int()))
	randomSuffixBuilder.WriteString(".")
	randomSuffixBuilder.WriteString(format)

	return randomSuffixBuilder.String()
}

func generateObjectFullPath(domainPrefix string, id string, format string) string {
	objectName := strings.Builder{}

	prefix := generateObjectPrefix(domainPrefix, id)
	randomSuffix := generateRandomSuffixWithFormat(format)

	objectName.WriteString(prefix)
	objectName.WriteString(randomSuffix)

	return objectName.String()
}

func createPreauthRequest(requestType objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeEnum, objectName string) objectstorage.CreatePreauthenticatedRequestRequest {
	return objectstorage.CreatePreauthenticatedRequestRequest{
		NamespaceName: common.String(NAMESPACE),
		BucketName:    common.String(BUCKETNAME),
		CreatePreauthenticatedRequestDetails: objectstorage.CreatePreauthenticatedRequestDetails{
			Name:                common.String(strconv.Itoa(rand.Int())),
			AccessType:          requestType,
			TimeExpires:         &common.SDKTime{Time: time.Now().Add(3 * time.Minute)},
			BucketListingAction: objectstorage.PreauthenticatedRequestBucketListingActionDeny,
			ObjectName:          common.String(objectName),
		},
	}
}

func CreateUploadObjectURI(domainPrefix string, id string, format string) (*model.Picture, error) {
	client := GetClient()
	objectName := generateObjectFullPath(domainPrefix, id, format)

	preAuthRequestRequest := createPreauthRequest(
		objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectwrite,
		objectName,
	)

	preAuth, err := client.CreatePreauthenticatedRequest(context.Background(), preAuthRequestRequest)
	if err != nil {
		slog.Error("Error when creating PAR for object upload")
		slog.Error(err.Error())

		return nil, err
	}

	picture := &model.Picture{
		Uri: *preAuth.FullPath,
	}

	return picture, nil
}

func ListObjectsByID(domainPrefix string, id string) (*model.Pictures, error) {
	slog.Info("OCI client: listing objects by id", "domainPrefix", domainPrefix, "id", id)

	objectPrefix := generateObjectPrefix(domainPrefix, id)
	listObjectRequest := objectstorage.ListObjectsRequest{
		NamespaceName: common.String(NAMESPACE),
		BucketName:    common.String(BUCKETNAME),
		Prefix:        common.String(objectPrefix),
	}

	slog.Info("Generated object prefix", "prefix", objectPrefix)

	response, err := obsClient.ListObjects(context.Background(), listObjectRequest)
	if err != nil {
		slog.Error("Cannot list objects", "listObjectRequest", listObjectRequest)
		slog.Error(err.Error())

		return nil, err
	}

	pictures := &model.Pictures{}
	for _, obj := range response.Objects {
		preAuthRequestRequest := createPreauthRequest(
			objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectread,
			*obj.Name,
		)

		slog.Info("Got object", "name", *obj.Name)

		response, err := obsClient.CreatePreauthenticatedRequest(context.Background(), preAuthRequestRequest)
		if err != nil {
			slog.Error("Cannot create PAR for listing objects")
			slog.Error(err.Error())

			return nil, err
		}

		slog.Info("Got response from OCI client", "fullpath", *response.FullPath)

		*pictures = append(*pictures, model.Picture{Uri: *response.FullPath})
	}

	return pictures, nil
}
