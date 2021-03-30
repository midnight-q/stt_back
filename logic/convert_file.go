package logic

import (
    "stt_back/common"
    "stt_back/errors"
    "stt_back/types"

    "github.com/jinzhu/gorm"
)

func ConvertFileFind(filter types.ConvertFileFilter)  (result []types.ConvertFile, totalRecords int, err error) {

	//ConvertFile Find logic code

    return
}

func ConvertFileMultiCreate(filter types.ConvertFileFilter)  (data []types.ConvertFile, err error) {

	//ConvertFile MultiCreate logic code

    return
}

func ConvertFileCreate(filter types.ConvertFileFilter, query *gorm.DB)  (data types.ConvertFile, err error) {
    model := filter.GetConvertFileModel()
    fileFormat := common.GetFileFormatFromName(model.InputFilename)

    resultFile := []byte{}
    switch fileFormat {
    case "wav":
        break
    case "mp3":
        break
    default:
        return types.ConvertFile{}, getUnsupportedFileFormatError()
    }

    // TODO: pass result file into external service

    // TODO: save converted input file on disk

    // TODO: save ConvertLog in db
	//ConvertFile Create logic code
    return
}

func getUnsupportedFileFormatError() error {
    return errors.NewErrorWithCode("Unsupported file format", errors.ErrorCodeUnsupportedFileFormat, "InputFilename")
}

func ConvertFileRead(filter types.ConvertFileFilter)  (data types.ConvertFile, err error) {

	//ConvertFile Read logic code
    return
}


func ConvertFileMultiUpdate(filter types.ConvertFileFilter)  (data []types.ConvertFile, err error) {

	//ConvertFile MultiUpdate logic code
    return
}

func ConvertFileUpdate(filter types.ConvertFileFilter, query *gorm.DB)  (data types.ConvertFile, err error) {

	//ConvertFile Update logic code
    return
}



func ConvertFileMultiDelete(filter types.ConvertFileFilter)  (isOk bool, err error) {

	//ConvertFile MultiDelete logic code
    return
}

func ConvertFileDelete(filter types.ConvertFileFilter, query *gorm.DB)  (isOk bool, err error) {

	//ConvertFile Delete logic code
    return
}



func ConvertFileFindOrCreate(filter types.ConvertFileFilter)  (data types.ConvertFile, err error) {
    
	//ConvertFile FindOrCreate logic code
    return
}


func ConvertFileUpdateOrCreate(filter types.ConvertFileFilter)  (data types.ConvertFile, err error) {
    
	//ConvertFile UpdateOrCreate logic code
    return
}

// add all assign functions
