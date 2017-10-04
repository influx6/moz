package ast

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/influx6/faux/metrics"
	"github.com/influx6/gobuild/build"
	"github.com/influx6/moz/gen"
)

var (
	// ErrEmptyList defines a error returned for a empty array or slice.
	ErrEmptyList = errors.New("Slice/List is empty")
)

// FilteredPackageWithBuildCtx parses the package directory which generates a series of ast with associated
// annotation for processing by using the golang token parser, it uses the build.Context to
// collected context details for the package and only processes the files found by the build context.
// If you need something more broad without filtering, use PackageWithBuildCtx.
func FilteredPackageWithBuildCtx(log metrics.Metrics, dir string, ctx build.Context) (RootPackage, error) {
	rootbuildPkg, err := ctx.ImportDir(dir, 0)
	if err != nil {
		log.Emit(metrics.Errorf("Failed to retrieve build.Package for root directory").
			With("file", dir).
			With("dir", dir).
			With("error", err.Error()).
			With("mode", build.FindOnly))
		return RootPackage{}, err
	}

	log.Emit(metrics.Info("Generated build.Package").
		With("file", dir).
		With("dir", dir).
		With("pkg", rootbuildPkg).
		With("mode", build.FindOnly))

	allowed := make(map[string]bool)
	for _, file := range rootbuildPkg.GoFiles {
		allowed[file] = true
	}

	filter := func(f os.FileInfo) bool {
		log.Emit(metrics.Info("Parse Filtering file").With("incoming-file", f.Name()).With("allowed", allowed[f.Name()]))
		return allowed[f.Name()]
	}

	var rpkg RootPackage
	rpkg.Path = rootbuildPkg.Dir
	rpkg.BuildPkg = rootbuildPkg
	rpkg.Package = rootbuildPkg.ImportPath

	tokenFiles := token.NewFileSet()
	packages, err := parser.ParseDir(tokenFiles, dir, filter, parser.ParseComments)
	if err != nil {
		log.Emit(metrics.Error(err).With("message", "Failed to parse dir").With("dir", dir))
		return rpkg, err
	}

	packageDeclrs := make(map[string]Package)
	packageBuilds := make(map[string]*build.Package)

	for _, pkg := range packages {
		var pkgFiles []string

		for path := range pkg.Files {
			pkgFiles = append(pkgFiles, path)
		}

		for path, file := range pkg.Files {
			pathPkg := filepath.Dir(path)
			buildPkg, ok := packageBuilds[pathPkg]
			if !ok {
				buildPkg, err = ctx.ImportDir(pathPkg, 0)
				if err != nil {
					log.Emit(metrics.Errorf("Failed to retrieve build.Package").
						With("file", path).
						With("dir", dir).
						With("file-dir", filepath.Dir(path)).
						With("error", err.Error()).
						With("mode", build.FindOnly))
					return RootPackage{}, err
				}

				packageBuilds[pathPkg] = buildPkg

				log.Emit(metrics.Info("Generated build.Package").
					With("file", path).
					With("pkg", buildPkg).
					With("file-dir", filepath.Dir(path)).
					With("dir", dir).
					With("mode", build.FindOnly))
			}

			res, err := parseFileToPackage(log, dir, path, pkg.Name, tokenFiles, file, pkg)
			if err != nil {
				log.Emit(metrics.Error(err).With("message", "Failed to parse file").With("dir", dir).With("file", file.Name.Name).With("Package", pkg.Name))
				return RootPackage{}, err
			}

			log.Emit(metrics.Info("Parsed Package File").With("dir", dir).With("file", file.Name.Name).With("path", path).With("Package", pkg.Name))

			if owner, ok := packageDeclrs[res.Package]; ok {
				owner.Packages = append(owner.Packages, res)
				packageDeclrs[res.Package] = owner
				continue
			}

			packageDeclrs[res.Package] = Package{
				Path:     res.Path,
				Package:  res.Package,
				BuildPkg: buildPkg,
				Files:    pkgFiles,
				Packages: []PackageDeclaration{res},
				Doc:      doc.New(pkg, buildPkg.ImportPath, doc.AllMethods),
			}
		}
	}

	rpkg.Packages = packageDeclrs

	return rpkg, nil
}

// PackageWithBuildCtx parses the package directory which generates a series of ast with associated
// annotation for processing by using the golang token parser, it uses the build.Context to
// collected context details for the package but does not use it has a means to select the files to
// process. PackageWithBuildCtx processes all files in package directory. If you want one which takes
// into consideration build.Context fields using FilteredPackageWithBuildCtx.
func PackageWithBuildCtx(log metrics.Metrics, dir string, ctx build.Context) (RootPackage, error) {
	rootbuildPkg, err := ctx.ImportDir(dir, 0)
	if err != nil {
		log.Emit(metrics.Errorf("Failed to retrieve build.Package for root directory").
			With("file", dir).
			With("dir", dir).
			With("error", err.Error()).
			With("mode", build.FindOnly))
		return RootPackage{}, err
	}

	log.Emit(metrics.Info("Generated build.Package").
		With("file", dir).
		With("dir", dir).
		With("pkg", rootbuildPkg).
		With("mode", build.FindOnly))

	tokenFiles := token.NewFileSet()
	packages, err := parser.ParseDir(tokenFiles, dir, nil, parser.ParseComments)
	if err != nil {
		log.Emit(metrics.Error(err).With("message", "Failed to parse directory").With("dir", dir))
		return RootPackage{}, err
	}

	packageDeclrs := make(map[string]Package)
	packageBuilds := make(map[string]*build.Package)

	for _, pkg := range packages {
		var pkgFiles []string

		for path := range pkg.Files {
			pkgFiles = append(pkgFiles, path)
		}

		for path, file := range pkg.Files {
			pathPkg := filepath.Dir(path)
			buildPkg, ok := packageBuilds[pathPkg]
			if !ok {
				buildPkg, err = ctx.ImportDir(pathPkg, 0)
				if err != nil {
					log.Emit(metrics.Errorf("Failed to retrieve build.Package").
						With("file", path).
						With("dir", dir).
						With("file-dir", filepath.Dir(path)).
						With("error", err.Error()).
						With("mode", build.FindOnly))
					return RootPackage{}, err
				}

				packageBuilds[pathPkg] = buildPkg

				log.Emit(metrics.Info("Generated build.Package").
					With("file", path).
					With("pkg", buildPkg).
					With("file-dir", filepath.Dir(path)).
					With("dir", dir).
					With("mode", build.FindOnly))
			}

			res, err := parseFileToPackage(log, dir, path, pkg.Name, tokenFiles, file, pkg)
			if err != nil {
				log.Emit(metrics.Error(err).With("message", "Failed to parse file").With("dir", dir).With("file", file.Name.Name).With("Package", pkg.Name))
				return RootPackage{}, err
			}

			log.Emit(metrics.Info("Parsed Package File").With("dir", dir).With("file", file.Name.Name).With("path", path).With("Package", pkg.Name))

			if owner, ok := packageDeclrs[res.Package]; ok {
				owner.Packages = append(owner.Packages, res)
				packageDeclrs[res.Package] = owner
				continue
			}

			packageDeclrs[res.Package] = Package{
				Path:     res.Path,
				Files:    pkgFiles,
				Package:  res.Package,
				BuildPkg: buildPkg,
				Packages: []PackageDeclaration{res},
				Doc:      doc.New(pkg, buildPkg.ImportPath, doc.AllMethods),
			}
		}
	}

	return RootPackage{
		Path:     rootbuildPkg.Dir,
		Package:  rootbuildPkg.ImportPath,
		BuildPkg: rootbuildPkg,
		Packages: packageDeclrs,
	}, nil
}

// PackageFileWithBuildCtx parses the package from the provided file.
func PackageFileWithBuildCtx(log metrics.Metrics, path string, ctx build.Context) (Package, error) {
	dir := filepath.Dir(path)
	fName := filepath.Base(path)

	buildPkg, err := ctx.ImportDir(dir, 0)
	if err != nil {
		log.Emit(metrics.Errorf("Failed to retrieve build.Package").
			With("file", path).
			With("dir", dir).
			With("error", err.Error()).
			With("mode", build.FindOnly))
		return Package{}, err
	}

	log.Emit(metrics.Info("Generated build.Package").
		With("file", path).
		With("dir", dir).
		With("pkg", buildPkg).
		With("mode", build.FindOnly))

	allowed := map[string]bool{
		fName: true,
	}

	filter := func(f os.FileInfo) bool {
		log.Emit(metrics.Info("Parse Filtering file").With("incoming-file", f.Name()).With("allowed", allowed[f.Name()]))
		return allowed[f.Name()]
	}

	tokenFiles := token.NewFileSet()
	packages, err := parser.ParseDir(tokenFiles, path, filter, parser.ParseComments)
	if err != nil {
		log.Emit(metrics.Error(err).With("message", "Failed to parse file").With("dir", dir).With("file", path))
		return Package{}, err
	}

	var pkg *ast.Package
	for _, pkg = range packages {
		if pkg.Name != buildPkg.Name {
			continue
		}
		break
	}

	var declrs []PackageDeclaration
	var pkgFiles []string

	for fpath, file := range pkg.Files {
		if fpath != path {
			continue
		}

		pkgFiles = append(pkgFiles, fpath)

		res, err := parseFileToPackage(log, dir, path, buildPkg.Name, tokenFiles, file, pkg)
		if err != nil {
			log.Emit(metrics.Error(err).With("message", "Failed to parse file").With("dir", dir).With("file", file.Name.Name).With("Package", pkg.Name))
			return Package{}, err
		}

		declrs = append(declrs, res)
	}

	return Package{
		Path:     buildPkg.Dir,
		Files:    pkgFiles,
		BuildPkg: buildPkg,
		Package:  buildPkg.ImportPath,
		Packages: declrs,
		Doc:      doc.New(pkg, buildPkg.ImportPath, doc.AllMethods),
	}, nil
}

//===========================================================================================================

// ParseFileAnnotations parses the package from the provided file.
func ParseFileAnnotations(log metrics.Metrics, path string) (Package, error) {
	return PackageFileWithBuildCtx(log, path, build.Default)
}

// ParseAnnotations parses the package which generates a series of ast with associated
// annotation for processing.
func ParseAnnotations(log metrics.Metrics, dir string) (RootPackage, error) {
	return PackageWithBuildCtx(log, dir, build.Default)
}

func parseFileToPackage(log metrics.Metrics, dir string, path string, pkgName string, tokenFiles *token.FileSet, file *ast.File, pkgAstObj *ast.Package) (PackageDeclaration, error) {
	var packageDeclr PackageDeclaration

	{
		pkgSource, _ := readSource(path)

		packageDeclr.Package = pkgName
		packageDeclr.FilePath = path
		packageDeclr.Source = string(pkgSource)
		packageDeclr.Imports = make(map[string]ImportDeclaration, 0)
		packageDeclr.ObjectFunc = make(map[*ast.Object][]FuncDeclaration, 0)

		for _, comment := range file.Comments {
			packageDeclr.Comments = append(packageDeclr.Comments, comment.Text())
		}

		for _, imp := range file.Imports {
			beginPosition, endPosition := tokenFiles.Position(imp.Pos()), tokenFiles.Position(imp.End())
			positionLength := endPosition.Offset - beginPosition.Offset
			source, err := readSourceIn(beginPosition.Filename, int64(beginPosition.Offset), positionLength)
			if err != nil {
				return packageDeclr, err
			}

			var pkgName string

			if imp.Name != nil {
				pkgName = imp.Name.Name
			}

			impPkgPath, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				impPkgPath = imp.Path.Value
			}

			var comment string

			if imp.Comment != nil {
				comment = imp.Comment.Text()
			}

			packageDeclr.Imports[pkgName] = ImportDeclaration{
				Comments: comment,
				Name:     pkgName,
				Path:     impPkgPath,
				Source:   string(source),
			}
		}

		if relPath, err := filepath.Rel(GoSrcPath, path); err == nil {
			packageDeclr.Path = filepath.Dir(relPath)
			packageDeclr.File = filepath.Base(relPath)
		}

		if runtime.GOOS == "windows" {
			packageDeclr.Path = filepath.ToSlash(packageDeclr.Path)
			packageDeclr.File = filepath.ToSlash(packageDeclr.File)
			packageDeclr.FilePath = filepath.ToSlash(packageDeclr.FilePath)
		}

		if file.Doc != nil {
			annotationRead := ReadAnnotationsFromCommentry(bytes.NewBufferString(file.Doc.Text()))

			log.Emit(metrics.Info("Annotations in Package comments").
				With("dir", dir).
				With("annotations", len(annotationRead)).
				With("file", file.Name.Name))

			packageDeclr.Annotations = append(packageDeclr.Annotations, annotationRead...)
		}

		// Collect and categorize annotations in types and their fields.
	declrLoop:
		for _, declr := range file.Decls {
			beginPosition, endPosition := tokenFiles.Position(declr.Pos()), tokenFiles.Position(declr.End())
			positionLength := endPosition.Offset - beginPosition.Offset
			source, err := readSourceIn(beginPosition.Filename, int64(beginPosition.Offset), positionLength)
			if err != nil {
				return packageDeclr, err
			}

			switch rdeclr := declr.(type) {
			case *ast.FuncDecl:
				var comment string

				if rdeclr.Doc != nil {
					comment = rdeclr.Doc.Text()
				}

				var annotations []AnnotationDeclaration
				associations := make(map[string]AnnotationAssociationDeclaration, 0)

				if rdeclr.Doc != nil {
					annotationRead := ReadAnnotationsFromCommentry(bytes.NewBufferString(rdeclr.Doc.Text()))

					for _, item := range annotationRead {
						log.Emit(metrics.Info("Annotation in Function Decleration comment").
							With("dir", dir).
							With("annotation", item.Name).
							With("position", rdeclr.Pos()))

						switch item.Name {
						case "associates":
							log.Emit(metrics.Error(errors.New("Association Annotation in Decleration is incomplete: Expects 3 elements")).
								With("dir", dir).
								With("association", item.Arguments).
								With("position", rdeclr.Pos()))

							if len(item.Arguments) >= 3 {
								associations[item.Arguments[0]] = AnnotationAssociationDeclaration{
									Record:     item,
									Template:   item.Template,
									Action:     item.Arguments[1],
									TypeName:   item.Arguments[2],
									Annotation: strings.TrimPrefix(item.Arguments[0], "@"),
								}
							}
						default:
							annotations = append(annotations, item)
						}
					}
				}

				var defFunc FuncDeclaration

				defFunc.Comments = comment
				defFunc.Source = string(source)
				defFunc.FuncDeclr = rdeclr
				defFunc.Type = rdeclr.Type
				defFunc.Position = rdeclr.Pos()
				defFunc.Path = packageDeclr.Path
				defFunc.File = packageDeclr.File
				defFunc.PackageDeclr = &packageDeclr
				defFunc.FuncName = rdeclr.Name.Name
				defFunc.Length = positionLength
				defFunc.From = beginPosition.Offset
				defFunc.Package = packageDeclr.Package
				defFunc.FilePath = packageDeclr.FilePath
				defFunc.Annotations = annotations
				defFunc.Associations = associations

				if rdeclr.Type != nil {
					defFunc.Returns = rdeclr.Type.Results
					defFunc.Arguments = rdeclr.Type.Params
				}

				if rdeclr.Recv != nil {
					defFunc.FuncType = rdeclr.Recv

					nameIdent := rdeclr.Recv.List[0]

					if receiverNameType, ok := nameIdent.Type.(*ast.Ident); ok {
						defFunc.RecieverName = receiverNameType.Name
						defFunc.Reciever = receiverNameType.Obj
						defFunc.RecieverIdent = receiverNameType

						if rems, ok := packageDeclr.ObjectFunc[receiverNameType.Obj]; ok {
							rems = append(rems, defFunc)
							packageDeclr.ObjectFunc[receiverNameType.Obj] = rems
						} else {
							packageDeclr.ObjectFunc[receiverNameType.Obj] = []FuncDeclaration{defFunc}
						}

						continue declrLoop
					}
				}

				packageDeclr.Functions = append(packageDeclr.Functions, defFunc)
				continue declrLoop

			case *ast.GenDecl:
				var comment string

				if rdeclr.Doc != nil {
					comment = rdeclr.Doc.Text()
				}

				var annotations []AnnotationDeclaration

				associations := make(map[string]AnnotationAssociationDeclaration, 0)

				if rdeclr.Doc != nil {
					annotationRead := ReadAnnotationsFromCommentry(bytes.NewBufferString(rdeclr.Doc.Text()))

					for _, item := range annotationRead {
						log.Emit(metrics.Info("Annotation in Decleration comment").
							With("dir", dir).
							With("annotation", item.Name).
							With("position", rdeclr.Pos()).
							With("token", rdeclr.Tok.String()))

						switch item.Name {
						case "associates":
							log.Emit(metrics.Error(errors.New("Association Annotation in Decleration is incomplete: Expects 3 elements")).
								With("dir", dir).
								With("association", item.Arguments).
								With("position", rdeclr.Pos()).
								With("token", rdeclr.Tok.String()))

							if len(item.Arguments) >= 3 {
								associations[item.Arguments[0]] = AnnotationAssociationDeclaration{
									Record:     item,
									Template:   item.Template,
									Action:     item.Arguments[1],
									TypeName:   item.Arguments[2],
									Annotation: strings.TrimPrefix(item.Arguments[0], "@"),
								}
							}
						default:
							annotations = append(annotations, item)
						}
					}
				}

				for _, spec := range rdeclr.Specs {
					switch obj := spec.(type) {
					case *ast.ValueSpec:
						// Handles variable declaration
						// i.e Spec:
						// &ast.ValueSpec{Doc:(*ast.CommentGroup)(nil), Names:[]*ast.Ident{(*ast.Ident)(0xc4200e4a00)}, Type:ast.Expr(nil), Values:[]ast.Expr{(*ast.BasicLit)(0xc4200e4a20)}, Comment:(*ast.CommentGroup)(nil)}
						// &ast.ValueSpec{Doc:(*ast.CommentGroup)(nil), Names:[]*ast.Ident{(*ast.Ident)(0xc4200e4a40)}, Type:(*ast.Ident)(0xc4200e4a60), Values:[]ast.Expr(nil), Comment:(*ast.CommentGroup)(nil)}

					case *ast.TypeSpec:

						switch robj := obj.Type.(type) {
						case *ast.StructType:

							log.Emit(metrics.Info("Annotation in Decleration").
								With("Type", "Struct").
								With("Annotations", len(annotations)).
								With("StructName", obj.Name))

							packageDeclr.Structs = append(packageDeclr.Structs, StructDeclaration{
								Object:       obj,
								Struct:       robj,
								Annotations:  annotations,
								Associations: associations,
								Source:       string(source),
								Comments:     comment,
								File:         packageDeclr.File,
								Package:      packageDeclr.Package,
								Path:         packageDeclr.Path,
								FilePath:     packageDeclr.FilePath,
								From:         beginPosition.Offset,
								Length:       positionLength,
								PackageDeclr: &packageDeclr,
							})
							break

						case *ast.InterfaceType:
							log.Emit(metrics.Info("Annotation in Decleration").
								With("Type", "Interface").
								With("Annotations", len(annotations)).
								With("StructName", obj.Name))

							packageDeclr.Interfaces = append(packageDeclr.Interfaces, InterfaceDeclaration{
								Object:       obj,
								Interface:    robj,
								Comments:     comment,
								Annotations:  annotations,
								Associations: associations,
								Source:       string(source),
								PackageDeclr: &packageDeclr,
								File:         packageDeclr.File,
								Package:      packageDeclr.Package,
								Path:         packageDeclr.Path,
								FilePath:     packageDeclr.FilePath,
								From:         beginPosition.Offset,
								Length:       positionLength,
							})
							break

						default:
							log.Emit(metrics.Info("Annotation in Decleration").
								With("Type", "OtherType").
								With("Marker", "NonStruct/NonInterface:Type").
								With("Annotations", len(annotations)).
								With("StructName", obj.Name))

							packageDeclr.Types = append(packageDeclr.Types, TypeDeclaration{
								Object:       obj,
								Annotations:  annotations,
								Comments:     comment,
								Associations: associations,
								Source:       string(source),
								PackageDeclr: &packageDeclr,
								File:         packageDeclr.File,
								Package:      packageDeclr.Package,
								Path:         packageDeclr.Path,
								FilePath:     packageDeclr.FilePath,
								From:         beginPosition.Offset,
								Length:       positionLength,
							})
						}

					case *ast.ImportSpec:
						// Do Nothing.
					}
				}

			case *ast.BadDecl:
				// Do Nothing.
			}
		}

	}

	return packageDeclr, nil
}

//===========================================================================================================

// Parse takes the provided packages parsing all internals declarations with the appropriate generators suited to the type and annotations.
func Parse(toDir string, log metrics.Metrics, provider *AnnotationRegistry, doFileOverwrite bool, pkgDeclrs ...Package) error {
	for _, pkg := range pkgDeclrs {
		if err := ParsePackage(toDir, log, provider, doFileOverwrite, pkg); err != nil {
			return err
		}
	}

	return nil
}

// WriteDirectives defines a function which houses the logic to write WriteDirective into file system.
func WriteDirectives(log metrics.Metrics, toDir string, doFileOverwrite bool, wds ...gen.WriteDirective) error {
	for _, wd := range wds {
		if err := WriteDirective(log, toDir, doFileOverwrite, wd); err != nil {
			return err
		}
	}

	return nil
}

// WriteDirective defines a function which houses the logic to write WriteDirective into file system.
func WriteDirective(log metrics.Metrics, toDir string, doFileOverwrite bool, item gen.WriteDirective) error {
	log.Emit(metrics.Info("Execute WriteDirective").With("File", item.FileName).With("Overwrite", item.DontOverride).With("Dir", item.Dir))

	if filepath.IsAbs(item.Dir) {
		err := errors.New("gen.WriteDirectiveError: Expected relative Dir path not absolute")
		log.Emit(metrics.Error(err).With("File", item.FileName).With("Overwrite", item.DontOverride).With("Dir", item.Dir))
		return err
	}

	namedFileDir := toDir
	if item.Dir != "" {
		namedFileDir = filepath.Join(toDir, item.Dir)
	}

	if err := os.MkdirAll(namedFileDir, 0700); err != nil && err != os.ErrExist {
		err = fmt.Errorf("IOError: Unable to create directory: %+q", err)
		log.Emit(metrics.Error(err).With("File", item.FileName).With("Overwrite", item.DontOverride).With("Dir", item.Dir))
		return err
	}

	if item.Writer == nil {
		log.Emit(metrics.Info("Resolved WriteDirective").With("File", item.FileName).With("Overwrite", item.DontOverride).With("Dir", item.Dir))
		return nil
	}

	if item.FileName == "" {
		err := fmt.Errorf("WriteDirective has no filename value attached")
		log.Emit(metrics.Error(err).With("File", item.FileName).With("Overwrite", item.DontOverride).With("Dir", item.Dir))
		return err
	}

	namedFile := filepath.Join(namedFileDir, item.FileName)

	fileStat, err := os.Stat(namedFile)
	if err == nil && !fileStat.IsDir() && item.DontOverride && !doFileOverwrite {
		log.Emit(metrics.Error(err).With("File", item.FileName).With("Overwrite", item.DontOverride).With("Dir", item.Dir).
			With("DestinationDir", namedFileDir).
			With("DestinationFile", namedFile))
		return err
	}

	newFile, err := os.Create(namedFile)
	if err != nil {
		log.Emit(metrics.Error(err).With("File", item.FileName).With("Overwrite", item.DontOverride).With("Dir", item.Dir).
			With("DestinationDir", namedFileDir).
			With("DestinationFile", namedFile))
		return err
	}

	defer newFile.Close()

	written, err := item.Writer.WriteTo(newFile)
	if err != nil && err != io.EOF {
		err = fmt.Errorf("IOError: Unable to write content to file: %+q", err)
		log.Emit(metrics.Error(err).With("File", item.FileName).With("Overwrite", item.DontOverride).With("Dir", item.Dir).
			With("DestinationDir", namedFileDir).
			With("DestinationFile", namedFile))
		return err
	}

	log.Emit(metrics.Info("Resolved WriteDirective").With("directive", item.DontOverride).
		With("data_written", written).
		With("DestinationDir", namedFileDir).
		With("DestinationFile", namedFile))

	return nil
}

// ParsePackage takes the provided package declrations parsing all internals with the appropriate generators suited to the type and annotations.
func ParsePackage(toDir string, log metrics.Metrics, provider *AnnotationRegistry, doFileOverwrite bool, pkgDeclrs Package) error {
	log.Emit(metrics.Info("Begin ParsePackage").With("toDir", toDir).
		With("overwriter-file", doFileOverwrite).
		With("package", pkgDeclrs.Package).
		With("doc", pkgDeclrs.Doc).
		With("doc.vars", len(pkgDeclrs.Doc.Vars)).
		With("doc.consts", len(pkgDeclrs.Doc.Consts)).
		With("doc.types", len(pkgDeclrs.Doc.Types)).
		With("doc.functions", len(pkgDeclrs.Doc.Funcs)))

	for _, pkg := range pkgDeclrs.Packages {
		log.Emit(metrics.Info("ParsePackage: Parse PackageDeclaration").
			With("toDir", toDir).With("overwriter-file", doFileOverwrite).
			With("package", pkg.Package).
			With("From", pkg.FilePath))

		wdrs, err := provider.ParseDeclr(pkgDeclrs, pkg, toDir)
		if err != nil {
			log.Emit(metrics.Error(fmt.Errorf("ParseFailure: Package %q", pkg.Package)).
				With("error", err.Error()).With("package", pkg.Package))
			continue
		}

		log.Emit(metrics.Info("ParseSuccess").With("From", pkg.FilePath).With("package", pkg.Package).With("Directives", len(wdrs)))

		for _, wd := range wdrs {
			if err := WriteDirective(log, toDir, doFileOverwrite, wd.WriteDirective); err != nil {
				log.Emit(metrics.Info("Annotation Resolved").With("annotation", wd.Annotation).
					With("dir", toDir).
					With("package", pkg.Package).
					With("file", pkg.File))
				continue
			}

			log.Emit(metrics.Info("Annotation Resolved").With("annotation", wd.Annotation).
				With("dir", toDir).
				With("package", pkg.Package).
				With("file", pkg.File))
		}

	}

	return nil
}

//===========================================================================================================

// WhichPackage is an utility function which returns the appropriate package name to use
// if a toDir is provided as destination.
func WhichPackage(toDir string, pkg Package) string {
	if toDir != "" && toDir != "./" && toDir != "." {
		return strings.ToLower(filepath.Base(toDir))
	}

	return pkg.Package
}

//===========================================================================================================
