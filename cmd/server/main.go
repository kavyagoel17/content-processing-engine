package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"io"
	"github.com/spf13/viper"
	"net/http"
)

type Config struct{
	Server ServerConfig `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`

}

type ServerConfig struct{
	Port int `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
}

type DatabaseConfig struct{
	URL string `mapstructure:"url"`
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Uploading File")
	file, handler, err := r.FormFile("myFile")
    if err != nil {
        ErrorLog.Println(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
	defer file.Close()

	InfoLog.Println(w, "Uploaded File: %s\n", handler.Filename)
    InfoLog.Println(w, "File Size: %d\n", handler.Size)
    InfoLog.Println(w, "MIME Header: %v\n", handler.Header)

	dst, err := os.Create(handler.Filename)
    if err != nil {
        http.Error(w, "Error saving the file", http.StatusInternalServerError)
        return
    }
	defer dst.Close()

}

func setupRoutes() {
    http.HandleFunc("/upload", uploadFile)
    http.ListenAndServe(":8080", nil)
}


func LoadConfig() (*Config, error){
	viper.SetDefault("server.port", 9090)
	viper.SetDefault("server.environment", "test")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err:= viper.ReadInConfig(); err!=nil{
		if _, ok := err.(viper.ConfigFileNotFoundError);
		!ok{
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		WarningLog.Println("Config file not found, relying on environment variables/defaults")
	}
	var config Config

	if err := viper.Unmarshal(&config); err!=nil{
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
var (
    WarningLog *log.Logger
    InfoLog   *log.Logger
    ErrorLog   *log.Logger
)

func init(){
	logFile, err := os.OpenFile("appLogs.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666)
	if err!=nil{
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)

	InfoLog = log.New(multiWriter, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(multiWriter, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(multiWriter, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)

}
func main(){

	cfg, err := LoadConfig()
	if err != nil {
		ErrorLog.Println("Failed to load configuration: %v", err)
	}
	InfoLog.Println("Loaded Server port:", cfg.Server.Port)
	InfoLog.Println("Loaded Server environment:", cfg.Server.Environment)
}