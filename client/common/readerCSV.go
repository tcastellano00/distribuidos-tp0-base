package common

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

func extractCSVFromZip(zipFilePath string, agencyID string) ([][]string, error) {
	// Abre el archivo ZIP
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Construye el nombre del archivo CSV esperado
	csvFileName := fmt.Sprintf("agency-%s.csv", agencyID)

	// Busca el archivo CSV en el archivo ZIP
	var csvFile *zip.File
	for _, file := range r.File {
		if file.Name == csvFileName {
			csvFile = file
			break
		}
	}

	if csvFile == nil {
		return nil, fmt.Errorf("archivo %s no encontrado en el ZIP", csvFileName)
	}

	// Lee el contenido del archivo CSV
	f, err := csvFile.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Lee los registros CSV
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func GetClienteMessagesFromZip(zipFilePath string, agencyID string) ([]*ClientMessage, error) {
	records, err := extractCSVFromZip(zipFilePath, agencyID)
	if err != nil {
		return nil, err
	}

	var clientMessages []*ClientMessage

	// Asume que el CSV tiene encabezados, por lo que empezamos desde la fila 1
	for _, record := range records {
		if len(record) < 5 {
			return nil, fmt.Errorf("registro CSV inválido: %v", record)
		}

		firstName := record[0]
		lastName := record[1]
		document := record[2]
		birthdate := record[3]
		number := record[4]

		clientMessage := CreateClientMessage(agencyID, firstName, lastName, document, birthdate, number)
		clientMessages = append(clientMessages, clientMessage)
	}

	return clientMessages, nil
}
