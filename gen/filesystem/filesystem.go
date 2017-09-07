package filesystem

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/influx6/moz/gen"
)

var (
	// bit flag for utf8 encoding
	utf8Encoding = uint16(1 << 11)
)

// MetaOption defines a function which receives a map to set values as desired.
type MetaOption func(map[string]interface{})

// MetaApply returns a MetaOption that applies a series of provided options to a map.
func MetaApply(ops ...MetaOption) MetaOption {
	return func(meta map[string]interface{}) {
		for _, op := range ops {
			op(meta)
		}
	}
}

// Meta sets the key(name) to the val in a map.
func Meta(name string, val interface{}) MetaOption {
	return func(meta map[string]interface{}) {
		meta[name] = val
	}
}

// MetaText sets the key(name) to the val in a map.
func MetaText(name string, message string, vals ...interface{}) MetaOption {
	return func(meta map[string]interface{}) {
		meta[name] = fmt.Sprintf(message, vals...)
	}
}

// Version sets the "version" key to the provided string in a map.
func Version(ver string) MetaOption {
	return func(meta map[string]interface{}) {
		meta["version"] = ver
	}
}

// Description sets the "description" key to the provided string in a map.
func Description(desc string, items ...interface{}) MetaOption {
	return func(meta map[string]interface{}) {
		meta["description"] = fmt.Sprintf(desc, items...)
	}
}

//================================================================================

// Content returns a io.WriterTo containing the provided data string.
func Content(data string) io.WriterTo {
	var bu bytes.Buffer
	bu.WriteString(data)
	return &bu
}

// ContentByte returns a io.WriterTo containing the provided data bytes.
func ContentByte(data []byte) io.WriterTo {
	var bu bytes.Buffer
	bu.Write(data)
	return &bu
}

// ContentFrom returns a io.WriterTo which wraps the io.Reader for piping
// the data from the giving reader.
func ContentFrom(r io.Reader) io.WriterTo {
	return &gen.FromReader{R: r}
}

//================================================================================

// FileWriter defines a structure that represent a file system file item
// with associated content.
type FileWriter struct {
	Name    string
	Content io.WriterTo
}

// File returns a instance of a FileWriter with associated name and content.
func File(name string, content io.WriterTo) FileWriter {
	var file FileWriter
	file.Name = name
	file.Content = content

	return file
}

// DirWriter implements a structure that represents a file system directory
// with associated files and name.
type DirWriter struct {
	Name       string
	ChildFiles []FileWriter
	ChildDirs  []DirWriter
}

// Dir returns a instance of a DirWriter with associated name.
// The files are scanned and are appropriately allocated if they are
// FileWriter or DirWriter.
func Dir(name string, files ...interface{}) DirWriter {
	var dir DirWriter
	dir.Name = name

	for _, item := range files {
		switch ritem := item.(type) {
		case FileWriter:
			dir.ChildFiles = append(dir.ChildFiles, ritem)
		case DirWriter:
			dir.ChildDirs = append(dir.ChildDirs, ritem)
		}
	}

	return dir
}

// GetDir returns the associated DirWriter for the giving relative filepath.
// Absolute path will be rejected.
func (dirs DirWriter) GetDir(dirPath string) (DirWriter, error) {
	if path.IsAbs(dirPath) {
		return DirWriter{}, errors.New("Absolute paths not allowed")
	}

	if dirPath == "" || dirPath == "." {
		return dirs, nil
	}

	levels := strings.Split(dirPath, "/")
	initial := levels[0]
	rest := path.Join(levels[1:]...)

	for _, dir := range dirs.ChildDirs {
		if dir.Name == initial {
			return dir.GetDir(rest)
		}
	}

	return DirWriter{}, fmt.Errorf("Dir %q not found in %q", initial, dirs.Name)
}

// GetFile returns the associated FileWriter for the giving relative filepath.
// Absolute path will be rejected.
func (dirs DirWriter) GetFile(filePath string) (FileWriter, error) {
	if path.IsAbs(filePath) {
		return FileWriter{}, errors.New("Absolute paths not allowed")
	}

	filedir, fileName := path.Split(filePath)
	if filedir == "" || filedir == "." {
		for _, file := range dirs.ChildFiles {
			if file.Name == fileName {
				return file, nil
			}
		}

		return FileWriter{}, fmt.Errorf("File %q not found in %q", fileName, dirs.Name)
	}

	levels := strings.Split(filedir, "/")
	initial := levels[0]
	rest := path.Join(append(levels[1:], fileName)...)

	for _, dir := range dirs.ChildDirs {
		if dir.Name == initial {
			return dir.GetFile(rest)
		}
	}

	return FileWriter{}, fmt.Errorf("File %q not found in %q", fileName, dirs.Name)
}

// Files runs through all child directory returning appropriate path and
// current item associated for that DirWriter.
func (dirs DirWriter) Files(rootDir string, cb func(hostFilePath string, hostFile FileWriter) error) error {
	for _, file := range dirs.ChildFiles {
		if err := cb(path.Join(rootDir, dirs.Name, file.Name), file); err != nil {
			return err
		}
	}
	return nil
}

// Dirs runs through all child directory returning appropriate path and
// current item associated for that DirWriter.
func (dirs DirWriter) Dirs(rootDir string, cb func(hostDirPath string, hostDir DirWriter) error) error {
	for _, dir := range dirs.ChildDirs {
		if err := cb(path.Join(rootDir, dirs.Name, dir.Name), dir); err != nil {
			return err
		}
	}
	return nil
}

// MemoryFileSystem defines a file system with associated files and directories which
// are to be converted into appropriate data stream by a consumer.
type MemoryFileSystem struct {
	Dir  DirWriter
	Meta map[string]interface{}
}

// FileSystem returns a instance of a FSWriter.
func FileSystem(content ...interface{}) MemoryFileSystem {
	var fsw MemoryFileSystem
	fsw.Meta = make(map[string]interface{})

	for _, item := range content {
		switch ritem := item.(type) {
		case func(map[string]interface{}):
			ritem(fsw.Meta)
		case MetaOption:
			ritem(fsw.Meta)
		case FileWriter:
			fsw.Dir.ChildFiles = append(fsw.Dir.ChildFiles, ritem)
		case DirWriter:
			fsw.Dir.ChildDirs = append(fsw.Dir.ChildDirs, ritem)
		}
	}

	return fsw
}

// GetDir returns the associated FileWriter for the giving relative filepath.
// Absolute path will be rejected.
func (mfs MemoryFileSystem) GetDir(dirPath string) (DirWriter, error) {
	return mfs.Dir.GetDir(dirPath)
}

// GetFile returns the associated FileWriter for the giving relative filepath.
// Absolute path will be rejected.
func (mfs MemoryFileSystem) GetFile(filePath string) (FileWriter, error) {
	return mfs.Dir.GetFile(filePath)
}

// Dirs runs through all filesystem child directories
// returning appropriate path and file to the provided callback.
func (mfs MemoryFileSystem) Dirs(cb func(hostFilePath string, hostFile DirWriter) error) error {
	return runThroughDirs(mfs.Dir, "", cb)
}

// Files runs through all filesystem files and child directory files
// returning appropriate path and file to the provided callback.
func (mfs MemoryFileSystem) Files(cb func(hostFilePath string, hostFile FileWriter) error) error {
	return runThroughFiles(mfs.Dir, "", cb)
}

// runThroughDirs runs all files within the DirWriter.
func runThroughDirs(base DirWriter, rootDir string, cb func(hostDirPath string, hostDir DirWriter) error) error {
	if err := cb(rootDir, base); err != nil {
		return err
	}

	return base.Dirs(rootDir, func(hostDirPath string, hostDir DirWriter) error {
		return runThroughDirs(hostDir, path.Dir(hostDirPath), cb)
	})
}

// runThroughFiles runs all files within the DirWriter.
func runThroughFiles(base DirWriter, rootDir string, cb func(hostDirPath string, hostDir FileWriter) error) error {
	if err := base.Files(rootDir, cb); err != nil {
		return err
	}

	return base.Dirs(rootDir, func(hostDirPath string, hostDir DirWriter) error {
		return runThroughFiles(hostDir, path.Dir(hostDirPath), cb)
	})
}

//======================================================================================

// JSONFileSystem implements io.WriteTo and transforms the MemoryFileSystem into a
// json hashmap using the encoding/json encoders.
type JSONFileSystem struct {
	FS     MemoryFileSystem
	indent bool
}

// JSONFS returns a new instance of the JSONFileSystem.
func JSONFS(fs MemoryFileSystem, indent bool) JSONFileSystem {
	return JSONFileSystem{FS: fs, indent: indent}
}

// ToReader returns a new reader from the contents of the GzipTarFileSystem.
// Each reader is unique and contains a complete data of all contents.
func (jfs *JSONFileSystem) ToReader() (io.Reader, error) {
	var comBuff bytes.Buffer
	if _, err := jfs.WriteTo(&comBuff); err != nil {
		return nil, err
	}

	return &comBuff, nil
}

// WriteTo implements io.WriterTo interface.
func (jfs *JSONFileSystem) WriteTo(w io.Writer) (int64, error) {
	archive := make(map[string]string)

	if metaJSON, err := json.Marshal(jfs.FS.Meta); err == nil {
		meta := File(".meta", ContentByte(metaJSON))
		_, err := handleJSONForFileWriter(archive, ".meta", meta)
		if err != nil {
			return 0, err
		}
	}

	if err := jfs.FS.Files(func(hostFilePath string, hostFile FileWriter) error {
		_, err := handleJSONForFileWriter(archive, hostFilePath, hostFile)
		return err
	}); err != nil {
		return 0, err
	}

	wc := gen.NewWriteCounter(w)
	encoder := json.NewEncoder(wc)

	if jfs.indent {
		encoder.SetIndent("\n", "\t")
	}

	if err := encoder.Encode(archive); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

func handleJSONForFileWriter(archive map[string]string, hostFilePath string, file FileWriter) (int64, error) {
	var bu bytes.Buffer

	total, err := file.Content.WriteTo(&bu)
	if err == nil {
		archive[hostFilePath] = bu.String()
	}

	return total, err
}

//======================================================================================

// ZipFileSystem implements io.WriteTo and transforms the MemoryFileSystem into a
// tar archive using the archive/zip writers.
type ZipFileSystem struct {
	FS MemoryFileSystem
}

// ZipFS returns a new instance of the ZipFileSystem.
func ZipFS(fs MemoryFileSystem) ZipFileSystem {
	return ZipFileSystem{FS: fs}
}

// ToReader returns a new reader from the contents of the GzipTarFileSystem.
// Each reader is unique and contains a complete data of all contents.
func (zfs *ZipFileSystem) ToReader() (io.Reader, error) {
	var comBuff bytes.Buffer
	if _, err := zfs.WriteTo(&comBuff); err != nil {
		return nil, err
	}

	return &comBuff, nil
}

// WriteTo implements io.WriterTo interface.
func (zfs *ZipFileSystem) WriteTo(w io.Writer) (int64, error) {
	archive := zip.NewWriter(w)
	defer archive.Close()

	var totalWritten int64

	if metaJSON, err := json.Marshal(zfs.FS.Meta); err == nil {
		meta := File(".meta", ContentByte(metaJSON))
		total, err := handleZipForFileWriter(archive, ".meta", meta)
		totalWritten += total

		if err != nil {
			return totalWritten, err
		}
	}

	if err := zfs.FS.Files(func(hostFilePath string, hostFile FileWriter) error {
		total, err := handleZipForFileWriter(archive, hostFilePath, hostFile)
		totalWritten += total

		return err
	}); err != nil {
		return 0, err
	}

	return totalWritten, nil
}

func handleZipForFileWriter(archive *zip.Writer, hostFilePath string, file FileWriter) (int64, error) {
	fileWriter, err := archive.CreateHeader(&zip.FileHeader{
		Name:   hostFilePath,
		Flags:  utf8Encoding,
		Method: zip.Deflate,
	})

	if err != nil {
		return 0, err
	}

	return file.Content.WriteTo(fileWriter)
}

//======================================================================================

// TarFileSystem implements io.WriteTo and transforms the MemoryFileSystem into a
// tar archive using the archive/tar writers.
type TarFileSystem struct {
	FS MemoryFileSystem
}

// TarFS returns a new instance of the TarFileSystem.
func TarFS(fs MemoryFileSystem) TarFileSystem {
	return TarFileSystem{FS: fs}
}

// ToReader returns a new reader from the contents of the GzipTarFileSystem.
// Each reader is unique and contains a complete data of all contents.
func (tfs *TarFileSystem) ToReader() (io.Reader, error) {
	var comBuff bytes.Buffer
	if _, err := tfs.WriteTo(&comBuff); err != nil {
		return nil, err
	}

	return &comBuff, nil
}

// WriteTo implements io.WriterTo interface.
func (tfs *TarFileSystem) WriteTo(w io.Writer) (int64, error) {
	archive := tar.NewWriter(w)
	defer archive.Close()

	var totalWritten int64

	if metaJSON, err := json.Marshal(tfs.FS.Meta); err == nil {
		meta := File(".meta", ContentByte(metaJSON))
		total, err := handleTarForFileWriter(archive, ".meta", meta)
		totalWritten += total

		if err != nil {
			return totalWritten, err
		}
	}

	if err := tfs.FS.Files(func(hostFilePath string, hostFile FileWriter) error {
		total, err := handleTarForFileWriter(archive, hostFilePath, hostFile)
		totalWritten += total
		return err
	}); err != nil {
		return 0, err
	}

	return totalWritten, nil
}

//======================================================================================

// GzipTarFileSystem implements io.WriteTo and transforms the MemoryFileSystem into a
// tar archive using the archive/tar writer to wrap a compress/gzip writer.
type GzipTarFileSystem struct {
	FS MemoryFileSystem
}

// GzipTarFS returns a new instance of the GzipTarFileSystem.
func GzipTarFS(fs MemoryFileSystem) GzipTarFileSystem {
	return GzipTarFileSystem{FS: fs}
}

// ToReader returns a new reader from the contents of the GzipTarFileSystem.
// Each reader is unique and contains a complete data of all contents.
func (gfs *GzipTarFileSystem) ToReader() (io.Reader, error) {
	var comBuff bytes.Buffer
	if _, err := gfs.WriteTo(&comBuff); err != nil {
		return nil, err
	}

	return &comBuff, nil
}

// WriteTo implements io.WriterTo interface.
func (gfs *GzipTarFileSystem) WriteTo(w io.Writer) (int64, error) {
	archive := tar.NewWriter(gzip.NewWriter(w))
	defer archive.Close()

	var totalWritten int64

	if metaJSON, err := json.Marshal(gfs.FS.Meta); err == nil {
		meta := File(".meta", ContentByte(metaJSON))
		total, err := handleTarForFileWriter(archive, ".meta", meta)
		totalWritten += total

		if err != nil {
			return totalWritten, err
		}
	}

	if err := gfs.FS.Files(func(hostFilePath string, hostFile FileWriter) error {
		total, err := handleTarForFileWriter(archive, hostFilePath, hostFile)
		totalWritten += total
		return err
	}); err != nil {
		return 0, err
	}

	return totalWritten, nil
}

//======================================================================================

func handleTarForFileWriter(archive *tar.Writer, hostFilePath string, file FileWriter) (int64, error) {
	var bu bytes.Buffer
	if n, err := file.Content.WriteTo(&bu); err != nil {
		return n, err
	}

	if err := archive.WriteHeader(&tar.Header{
		Name:    hostFilePath,
		Mode:    0666,
		ModTime: time.Now(),
		Size:    int64(bu.Len()),
	}); err != nil {
		return 0, err
	}

	return bu.WriteTo(archive)
}
