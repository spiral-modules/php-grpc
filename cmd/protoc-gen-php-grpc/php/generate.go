// MIT License
//
// Copyright (c) 2018 SpiralScout
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package php

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"strconv"
	"strings"
)

// Generate generates needed service classes
func Generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	resp := &plugin.CodeGeneratorResponse{}

	filePerMethod := false
	configName := "FilePerMethod="
	for _, config := range strings.Split(req.GetParameter(), ",") {
		if strings.HasPrefix(config, configName) {
			val, err :=  strconv.ParseBool(strings.TrimPrefix(config, configName))
			if err != nil {
				panic(err)
			}
			filePerMethod = val
		}
	}
	for _, file := range req.ProtoFile {
		for _, service := range file.Service {
			if filePerMethod == true {
				for _, method := range service.Method {
					resp.File = append(resp.File, generateMethod(req, file, service, method))
				}
			} else {
				resp.File = append(resp.File, generate(req, file, service))
			}
		}
	}

	return resp
}

func generate(
	req *plugin.CodeGeneratorRequest,
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
) *plugin.CodeGeneratorResponse_File {
	return &plugin.CodeGeneratorResponse_File{
		Name:    str(filename(file, service.Name)),
		Content: str(body(req, file, service)),
	}
}

func generateMethod(
	req *plugin.CodeGeneratorRequest,
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
	method *descriptor.MethodDescriptorProto,
) *plugin.CodeGeneratorResponse_File {
	return &plugin.CodeGeneratorResponse_File{
		Name:    str(methodInterfaceFilename(file, method.Name)),
		Content: str(methodInterfaceBody(req, file, service, method)),
	}
}

// helper to convert string into string pointer
func str(str string) *string {
	return &str
}
