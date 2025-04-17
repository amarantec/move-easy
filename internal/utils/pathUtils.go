package utils

import (
    "os"
    "github.com/amarantec/move-easy/internal"
    "log"
    "errors"
    "fmt"
    "net/http"
    "path/filepath"
    "html/template"
	"github.com/joho/godotenv"
)

var Templates *template.Template


func FindDir(rootPath string, targetName string) (string, error) {
	// Verificação direta para caminhos absolutos esperados em Docker
	if targetName == internal.TEMPLATES {
		defaultTemplatesPath := internal.DEFAULT_TEMPLATES_PATH
		if stat, err := os.Stat(defaultTemplatesPath); err == nil && stat.IsDir() {
			return defaultTemplatesPath, nil
		}
	}

	if targetName == internal.WWW {
		defaultWWWPath := internal.DEFAULT_WWW_PATH
		if stat, err := os.Stat(defaultWWWPath); err == nil && stat.IsDir() {
			return defaultWWWPath, nil
		}
	}

	// Busca recursiva padrão (fallback para ambiente local)
	files, err := os.ReadDir(rootPath)
	if err != nil {
		return internal.EMPTY, err
	}

	for _, file := range files {
		if file.IsDir() {
			if file.Name() == targetName {
				return filepath.Join(rootPath, file.Name()), nil
			}

			subDir, err := FindDir(filepath.Join(rootPath, file.Name()), targetName)
			if err == nil {
				return subDir, nil
			}
		}
	}

	return internal.EMPTY, errors.New("directory not found: " + targetName)
}

func LoadHTMLTemplates() {
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatal("error getting current working directory: ", err)
    }

    templatesDir, err := FindDir(cwd, internal.TEMPLATES)
    if err != nil {
        log.Fatal("templates directory not found: ", err)
    }

    pattern := filepath.Join(templatesDir, internal.HTML)
    Templates, err = template.ParseGlob(pattern)
    if err != nil {
        log.Fatalf("error loading templates: %v\n", err)
    }
}

func StaticFileHandler(urlPrefix string) (string, http.Handler) {
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatal("error getting working directory: ", err)
    }

    wwwDir, err := FindDir(cwd, internal.WWW)
    if err != nil {
        log.Fatal("www directory not found: ", err)
    }

    fs := http.FileServer(http.Dir(wwwDir))
    return urlPrefix, http.StripPrefix(urlPrefix, fs)
}

func findEnv(path string) (string, error) {
	filePath := filepath.Join(path, internal.ENV)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return internal.EMPTY, err
	}

	for _, file := range files {
		if file.IsDir() {
			subDirPath := filepath.Join(path, file.Name())
			subDirFilePath, err := findEnv(subDirPath)
			if err == nil {
				return subDirFilePath, nil
			}
		}
	}

	return internal.EMPTY, fmt.Errorf("file .env not found in %s or anywhere\n", path)
}

func LoadEnv() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting actual dir")
	}

	appPath := filepath.Dir(filepath.Dir(path))
	pathFile, err := findEnv(appPath)
	if err != nil {
		log.Fatal(".env no found")
	}
	godotenv.Load(pathFile)
}

func BuildConnectionString() (string, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return internal.EMPTY, fmt.Errorf("one or more environment variables are not set")
	}

	connectionString := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		dbHost, dbPort, dbUser, dbPassword, dbName)

	return connectionString, nil
}
