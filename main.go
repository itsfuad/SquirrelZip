package main

import (
	"file-compressor/compressor"
	"file-compressor/constants"
	"file-compressor/encryption"
	"file-compressor/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)


func handleDecompress(fileName, outputDir, password string) {
	encryptedFile, err := os.Open(fileName)
	if err != nil {
		utils.ColorPrint(utils.RED, fmt.Sprintf(constants.FILE_OPEN_ERROR, err.Error()))
		os.Exit(-1)
	}

	defer encryptedFile.Close()

	decryptedFilePath := fileName + ".decrypted"
	decryptedFile, err := os.Create(decryptedFilePath)
	if err != nil {
		utils.ColorPrint(utils.RED, fmt.Sprintf(constants.FILE_CREATE_ERROR, err.Error())+"\n")
		os.Exit(-1)
	}

	err = encryption.DecryptStream(encryptedFile, decryptedFile, password)
	if err != nil {
		utils.ColorPrint(utils.RED, fmt.Sprintf(constants.FAILED_TO_DECRYPT, err.Error())+"\n")
		//release file
		decryptedFile.Close()
		// delete the decrypted file
		utils.SafeDeleteFile(decryptedFilePath)
		os.Exit(-1)
	}

	decryptedFile.Close()

	paths, err := compressor.Decompress(decryptedFilePath, outputDir)
	if err != nil {
		utils.ColorPrint(utils.RED, err.Error()+"\n")
		// delete the decrypted file
		//utils.SafeDeleteFile(decryptedFilePath)
		os.Exit(-1)
	}

	// delete the decrypted file
	utils.SafeDeleteFile(decryptedFilePath)

	for _, path := range paths {
		utils.ColorPrint(utils.GREEN, "Output file: "+path+"\n")
	}
}

func handleCompress(fileNames []string, outputDir, password, algorithm string) {
	outputPath, fileMeta, err := compressor.Compress(fileNames, outputDir, algorithm)
	if err != nil {
		utils.ColorPrint(utils.RED, err.Error()+"\n")
		utils.SafeDeleteFile(outputPath)
		os.Exit(-1)
	}

	fileMeta.PrintFileInfo()
	fileMeta.PrintCompressionRatio()

	compressedFile, err := os.Open(outputPath)
	if err != nil {
		utils.ColorPrint(utils.RED, fmt.Sprintf(constants.FILE_OPEN_ERROR, err.Error())+"\n")
		os.Exit(-1)
	}

	fileName := outputPath
	fileExt := filepath.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, fileExt)

	finalFileName := utils.InvalidateFileName(fileName+".sq", "")

	finalFile, err := os.Create(finalFileName)
	if err != nil {
		utils.ColorPrint(utils.RED, fmt.Sprintf(constants.FILE_CREATE_ERROR, err.Error())+"\n")
		//release file
		compressedFile.Close()
		utils.SafeDeleteFile(outputPath)
		os.Exit(-1)
	}

	err = encryption.EncryptStream(compressedFile, finalFile, password)
	if err != nil {
		utils.ColorPrint(utils.RED, fmt.Sprintf(constants.FAILED_TO_ENCRYPT, err.Error())+"\n")
		//release file
		compressedFile.Close()
		finalFile.Close()
		utils.SafeDeleteFile(outputPath)
		utils.SafeDeleteFile(finalFileName)
		os.Exit(-1)
	}

	utils.ColorPrint(utils.GREEN, "Output file: "+finalFileName+"\n")

	compressedFile.Close()
	// delete the compressed file
	utils.SafeDeleteFile(outputPath)
}


func main() {

	startTime := time.Now()

	//cli arguments
	filenameStrs, outputDir, password, mode, algorithm := utils.ParseCLI()

	if mode == utils.DECOMPRESS {
		handleDecompress(filenameStrs[0], outputDir, password)
	} else {
		handleCompress(filenameStrs, outputDir, password, algorithm)
	}

	endTime := time.Now()
	utils.ColorPrint(utils.GREEN, "Time taken: "+utils.TimeTrack(startTime, endTime)+"\n")
}
