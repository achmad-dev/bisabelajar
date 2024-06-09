package firebasestore

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseClient holds the Firebase storage client
type FirebaseClient struct {
	storageClient *storage.BucketHandle
	bucketName    string
}

// NewFirebaseClient initializes a new Firebase storage client
func NewFirebaseClient(ctx context.Context, serviceAccountKeyPath string, bucketName string) (*FirebaseClient, error) {
	// Initialize the app with a service account, granting admin privileges
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	// Obtain a storage client from the app
	storageClient, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obtaining storage client: %v", err)
	}

	client, err := storageClient.DefaultBucket()
	if err != nil {
		return nil, fmt.Errorf("error getting default bucket: %v", err)
	}

	return &FirebaseClient{storageClient: client, bucketName: bucketName}, nil
}

// UploadFile uploads a file to Firebase storage
func (fc *FirebaseClient) UploadFile(ctx context.Context, localFilePath string, remoteFileName string) (string, error) {

	// Open the local file
	f, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening local file: %v", err)
	}
	defer f.Close()

	// Create a new object in the bucket
	wc := fc.storageClient.Object(remoteFileName).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	return fmt.Sprintf("gs://%s/%s", fc.bucketName, remoteFileName), nil
}

// DownloadFile downloads a file from Firebase storage
func (fc *FirebaseClient) DownloadFile(ctx context.Context, remoteFileName string, localFilePath string) error {

	// Create a new file locally
	f, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("error creating local file: %v", err)
	}
	defer f.Close()

	// Read the object from the bucket
	rc, err := fc.storageClient.Object(remoteFileName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("error reading object: %v", err)
	}
	defer rc.Close()

	// Copy the object's content to the local file
	if _, err = io.Copy(f, rc); err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}

	return nil
}

// DeleteFile deletes a file from Firebase storage
func (fc *FirebaseClient) DeleteFile(ctx context.Context, remoteFileName string) error {

	// Delete the object from the bucket
	err := fc.storageClient.Object(remoteFileName).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	return nil
}
