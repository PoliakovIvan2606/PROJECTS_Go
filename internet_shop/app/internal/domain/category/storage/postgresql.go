package storage

import (
	"internet_shop/pkg/client/postgresql"

	"github.com/sirupsen/logrus"
)

type storage struct {
	client postgresql.Client
	logger *logrus.Logger
}