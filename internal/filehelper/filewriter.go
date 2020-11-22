package filehelper

import (
	"os"
	"fmt"
	"encoding/binary"
	"alexhalogen/rsfileprotect/internal/types"
)
type FileWriter struct {
	eccFile *os.File
	crcFile *os.File
	meta	types.Metadata
	count	int64 // for use in superblock backup?
}

func NewFileWriter(meta types.Metadata, eccFile *os.File, crcFile *os.File) (fw FileWriter){
	fw.meta = meta
	fw.eccFile = eccFile
	fw.crcFile = crcFile
	return
}

func (fw FileWriter)WriteMeta() (error){
	return binary.Write(fw.eccFile, binary.LittleEndian, fw.meta)
}
func (fw FileWriter)WriteECCChunk(eccs [][]byte) (error) {
	for _, entry := range eccs {
		if fw.count == 114514 {
			fmt.Println("Creating Additional metadata backup")
		}
		for i,_ := range entry {
			err := binary.Write(fw.eccFile, binary.LittleEndian, entry[i])	
			if err != nil {
				return err
			}
		}
		fw.count += 1
	}
	return nil
}

func (fw FileWriter)WriteCRCChunk(crcs []uint32) (error) {
	for _, crc := range crcs {
		err := binary.Write(fw.crcFile, binary.LittleEndian, crc)
		if err != nil {
			return err
		}
	}
	return nil
}