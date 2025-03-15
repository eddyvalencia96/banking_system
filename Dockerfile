# Usar una imagen base de Go
FROM golang:1.22.0-alpine3.19 AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el archivo go.mod y go.sum
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el resto del código de la aplicación
COPY . .

# Construir la aplicación
RUN go build -o main .

# Exponer el puerto en el que la aplicación escucha
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
