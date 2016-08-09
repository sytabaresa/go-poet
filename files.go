package gopoet

import "bytes"

type FileSpec struct {
	Package                string
	InitializationPackages []Import
	Init                   CodeBlock
	CodeBlocks             []CodeBlock
}

func NewFileSpec(pkg string) *FileSpec {
	return &FileSpec{
		Package:                pkg,
		InitializationPackages: []Import{},
		CodeBlocks:             []CodeBlock{},
	}
}

func (f *FileSpec) String() string {
	var buffer bytes.Buffer
	didStartImportBlock := false

	buffer.WriteString("package " + f.Package + "\n\n")
	seen := map[string]struct{}{}

	var packages []Import
	for _, blk := range f.CodeBlocks {
		packages = append(packages, blk.GetImports()...)
	}

	for _, imp := range f.InitializationPackages {
		if _, found := seen[imp.GetPackage()]; !found && imp.GetPackage() != "" {
			if !didStartImportBlock {
				buffer.WriteString("import (\n")
				didStartImportBlock = true
			}

			buffer.WriteString("\t_ ")
			buffer.WriteString("\"" + imp.GetPackage() + "\"\n")
			seen[imp.GetPackage()] = struct{}{}
		}
	}

	for _, imp := range packages {
		if _, found := seen[imp.GetPackage()]; !found && imp.GetPackage() != "" {
			if !didStartImportBlock {
				buffer.WriteString("import (\n")
				didStartImportBlock = true
			}

			buffer.WriteString("\t")
			if imp.GetAlias() != "" {
				buffer.WriteString(imp.GetAlias())
				buffer.WriteString(" ")
			}
			buffer.WriteString("\"" + imp.GetPackage() + "\"\n")
			seen[imp.GetPackage()] = struct{}{}
		}
	}

	if didStartImportBlock {
		buffer.WriteString(")\n\n")
	}

	if f.Init != nil {
		f.CodeBlocks = append([]CodeBlock{f.Init}, f.CodeBlocks...)
	}

	for _, codeBlk := range f.CodeBlocks {
		buffer.WriteString(codeBlk.String())
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func (f *FileSpec) InitializationPackage(imp Import) *FileSpec {
	f.InitializationPackages = append(f.InitializationPackages, imp)
	return f
}

func (f *FileSpec) CodeBlock(blk CodeBlock) *FileSpec {
	f.CodeBlocks = append(f.CodeBlocks, blk)
	return f
}

func (f *FileSpec) InitFunction(blk CodeBlock) *FileSpec {
	f.Init = blk
	return f
}

func (f *FileSpec) GlobalVariable(name string, typ TypeReference, format string, args ...interface{}) *FileSpec {
	v := &Variable{
		Identifier: Identifier{
			Name: name,
			Type: typ,
		},
		Constant: false,
		Format:   format,
		Args:     args,
	}
	f.CodeBlocks = append(f.CodeBlocks, v)
	return f
}

func (f *FileSpec) GlobalConstant(name string, typ TypeReference, format string, args ...interface{}) *FileSpec {
	v := &Variable{
		Identifier: Identifier{
			Name: name,
			Type: typ,
		},
		Constant: true,
		Format:   format,
		Args:     args,
	}
	f.CodeBlocks = append(f.CodeBlocks, v)
	return f
}

func (f *FileSpec) VariableGrouping() *VariableGrouping {
	v := &VariableGrouping{}
	f.CodeBlocks = append(f.CodeBlocks, v)
	return v
}
