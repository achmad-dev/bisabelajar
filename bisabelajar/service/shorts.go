package service

import (
	middlewarev1 "bisabelajar/api/v1/middleware"
	"bisabelajar/dto"
	"context"
	"fmt"

	firebasestore "internal/pkg/firebase_store"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type ShortService struct {
	firebasestore *firebasestore.FirebaseClient
	tracer        trace.Tracer
	log           *logrus.Logger
}

func NewShortService(
	firebaseStoreClient *firebasestore.FirebaseClient,
	log *logrus.Logger,
	tracer trace.Tracer) ShortService {

	return ShortService{
		firebasestore: firebaseStoreClient,
		log:           log,
		tracer:        tracer}
}

func (b *ShortService) UploadShort(filePath string, shortDTO dto.ShortDTO, requestDetails middlewarev1.RequestDetails) error {
	remotefilename, err := b.firebasestore.UploadFile(context.Background(), filePath, fmt.Sprintf("%s%s", uuid.New(), filepath.Ext(filePath)))
	if err != nil {
		b.log.Errorf("error when upload with error %s", err.Error())
	}
	b.log.Infof("file uploaded in https://storage.googleapis.com/bookshello-300a4.appspot.com/%s", remotefilename)
	return nil
}
