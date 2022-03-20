/*
Copyright © 2022 K8sCommerce
*/
package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type categoryBuilder struct {
	apiURL    string
	storeKey  string
	outputDir string
	baseName  string
}

func NewCategoryBuilder(apiURL, storeKey, outputDir, baseName string) Builder {
	return &categoryBuilder{
		apiURL:    apiURL,
		storeKey:  storeKey,
		outputDir: outputDir,
		baseName:  baseName,
	}
}

func (b *categoryBuilder) Build() {
	categories := b.getCategories()
	for _, category := range categories {
		dir := path.Clean(fmt.Sprintf("%s/%s/%s", b.outputDir, b.baseName, category.Slug))
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "\nerror: %s\n\n", err.Error())
			os.Exit(1)
		}
		// write the file
		b.write(dir, "_index.md", category)
	}
}

func (b *categoryBuilder) write(dir, filename string, category Category) {
	mdFile := path.Clean(fmt.Sprintf("%s/%s", dir, filename))
	if _, err := os.Stat(mdFile); !errors.Is(err, os.ErrNotExist) {
		if err := os.Truncate(mdFile, 0); err != nil {
			log.Printf("Failed to truncate page: %s %v", mdFile, err)
		}
	}

	// Create a file for writing
	f, err := os.Create(mdFile)
	if err != nil {
		// failed to create/open the file
		fmt.Fprintf(os.Stderr, "\nerror: %s\n\n", err.Error())
		os.Exit(1)
	}
	f.WriteString("---\n")
	enc := yaml.NewEncoder(f)
	// enc.SetIndent("", "    ")
	if err := enc.Encode(category); err != nil {
		// failed to encode
		fmt.Fprintf(os.Stderr, "\nerror: %s\n\n", err.Error())
		os.Exit(1)
	}
	f.WriteString("---")
	if err := f.Close(); err != nil {
		// failed to close the file
		fmt.Fprintf(os.Stderr, "\nerror: %s\n\n", err.Error())
		os.Exit(1)
	}
}

func (b *categoryBuilder) getCategories() []Category {
	var categories []Category
	url := b.apiURL + "/v1/categories"

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nerror: %s\n\n", err.Error())
		os.Exit(1)
		return nil
	}
	req.Header.Set("Store-Key", b.storeKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nerror: Error when sending request to the server: %s\n\n", err.Error())
		os.Exit(1)
		return nil
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nerror: %s\n\n", err.Error())
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusOK {
		var getAllCategoriesResponse GetAllCategoriesResponse
		err = json.Unmarshal(responseBody, &getAllCategoriesResponse)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nerror: %s\n\n", err.Error())
			os.Exit(1)
		}

		categories = append(categories, getAllCategoriesResponse.Categories...)

	} else {
		fmt.Fprintf(os.Stderr, "\n\nerror: %s\n%s\n", resp.Status, string(responseBody))
		os.Exit(1)
	}
	return categories
}

type GetAllCategoriesResponse struct {
	Categories []Category `json:"categories" yaml:"categories"` // a collection of Category
}

type Category struct {
	Id              int64  `json:"id" yaml:"id"`                             // category id
	ParentId        int64  `json:"parent_id" yaml:"parent_id"`               // parent category id. references Category.Id
	Slug            string `json:"slug" yaml:"slug"`                         // slug name of the category
	Name            string `json:"name" yaml:"name"`                         // name of category
	Description     string `json:"description" yaml:"description"`           // description of category
	MetaTitle       string `json:"meta_title" yaml:"meta_title"`             // metatag title for SEO
	MetaDescription string `json:"meta_description" yaml:"meta_description"` // metatag description for SEO
	MetaKeywords    string `json:"meta_keywords" yaml:"meta_keywords"`       // metatag keywords for SEO
	Depth           int32  `json:"depth" yaml:"depth"`                       // category level depth
	SortOrder       int32  `json:"sort_order" yaml:"sort_order"`             // sort order of menu items on the same level and same parent id
}
