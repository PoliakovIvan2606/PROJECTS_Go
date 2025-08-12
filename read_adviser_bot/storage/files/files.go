package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"read_adviser_bot/lib/e"
	"read_adviser_bot/storage"
	"time"
)

const defaultPerm = 0774

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {e.Wrap("can't save page", err)}()
	// Формируем птуь до дериктории , где хранятся файлы
	fPath := filepath.Join(s.basePath, page.UserName)
	log.Print(fPath)
	
	// Создаём все нужные дериктории
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	// Формируем имя файла
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	// Дописываем имя файла пути
	fPath = filepath.Join(fPath, fName)
	
	// Создаём файл
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() {_=file.Close()}()
	
	// И записываем страницу в нужном формате
	if err:=gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() {e.Wrap("can't save page", err)}()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("con't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("con't remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if files %s exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	
	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil

	case err != nil:
		msg := fmt.Sprintf("can't check if files %s exists", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {_=f.Close()}()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode path", err)
	}

	return &p, nil
}


func fileName(p *storage.Page) (string, error) {
	return p.Hash()
} 