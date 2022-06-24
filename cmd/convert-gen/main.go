// Copyright © 2022 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

type Option struct {
	Src  string
	Dest string
}

func main() {
	opt := &Option{
		Src:  ".output",
		Dest: ".dest",
	}

	flag.StringVar(&opt.Src, "src", opt.Src, "specify the src parsing path")
	flag.StringVar(&opt.Dest, "dest", opt.Dest, "specify the dest convert path")
	flag.Parse()

	if err := run(opt); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Convert Finished")
}

func run(opt *Option) error {
	stat, err := os.Stat(opt.Src)
	if err != nil {
		return fmt.Errorf("stat path(%s) failed: %v", opt.Src, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("path(%s) is not dir", opt.Src)
	}

	dir, err := os.ReadDir(opt.Src)
	if err != nil {
		return fmt.Errorf("read dir(%s) failed: %v", opt.Src, err)
	}

	if err := os.MkdirAll(opt.Dest, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir(%s) failed: %v", opt.Dest, err)
	}
	for _, v := range dir {
		if v.IsDir() {
			continue
		}

		fullPath := filepath.Join(opt.Src, v.Name())
		filedata, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("read file(%s) failed: %v", fullPath, err)
		}

		crd := new(apiextensionsv1.CustomResourceDefinition)
		if err := yaml.Unmarshal(filedata, &crd); err != nil {
			return fmt.Errorf("unmarshal CustomResourceDefinition json failed: %v", err)
		}
		if len(crd.Spec.Versions) == 0 {
			continue
		}

		currentDir := filepath.Join(opt.Dest, crd.Spec.Group, crd.Spec.Names.Kind)
		if err := os.MkdirAll(currentDir, os.ModePerm); err != nil {
			return err
		}

		for _, version := range crd.Spec.Versions {
			crdCopy := crd.DeepCopy()
			crdCopy.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{version}
			b, err := yaml.Marshal(crdCopy)
			if err != nil {
				return fmt.Errorf("marshal CRD(%s) failed: %v", crdCopy.GroupVersionKind().String(), err)
			}

			currentFile := filepath.Join(currentDir, version.Name+".yaml")
			if err := os.WriteFile(currentFile, b, os.ModePerm); err != nil {
				return fmt.Errorf("write file(%s) failed: %v", currentFile, err)
			}
		}
	}
	return nil
}
