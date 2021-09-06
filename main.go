package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

	"muzzammil.xyz/jsonc"
)

type Theme struct {
	Name string
	Path string
	File string
	Dark bool
}

func main() {
	themes := []Theme{
		{
			Name: "dracula",
			Path: "dracula-theme.theme-dracula-2.22.3.vsix",
			File: "extension/theme/dracula.json",
			Dark: true,
		},
		{
			Name: "solarized-light",
			Path: "ryanolsonx.solarized-2.0.3.vsix",
			File: "extension/themes/light-color-theme.json",
			Dark: false,
		},
		{
			Name: "solarized-dark",
			Path: "ryanolsonx.solarized-2.0.3.vsix",
			File: "extension/themes/dark-color-theme.json",
			Dark: true,
		},
		{
			Name: "material-light",
			Path: "Equinusocio.vsc-material-theme-33.2.2.vsix",
			File: "extension/build/themes/Material-Theme-Lighter.json",
			Dark: false,
		},
		{
			Name: "material-dark",
			Path: "Equinusocio.vsc-material-theme-33.2.2.vsix",
			File: "extension/build/themes/Material-Theme-Default.json",
			Dark: true,
		},
		{
			Name: "github-light",
			Path: "GitHub.github-vscode-theme-3.0.0.vsix",
			File: "extension/themes/light.json",
			Dark: false,
		},
		{
			Name: "github-dark",
			Path: "GitHub.github-vscode-theme-3.0.0.vsix",
			File: "extension/themes/dark.json",
			Dark: true,
		},
	}

	for _, theme := range themes {
		fmt.Println("Process theme: ", theme.Name)

		fmt.Println("  Extract theme")
		content, err := extractTheme(theme)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("  Generate template")
		generateTheme(theme, content)
	}
}

func generateTheme(theme Theme, content []byte) {
	params := makeTemplateParams(theme, content)

	templateCss, err := template.ParseFiles("./theme-template.css")
	if err != nil {
		log.Fatal(err)
	}
	outCss, err := os.Create("./themes/" + theme.Name + ".css")
	if err != nil {
		log.Fatal(err)
	}
	err = templateCss.Execute(outCss, params)
	if err != nil {
		log.Fatal(err)
	}

	templateJs, err := template.ParseFiles("./theme-template.js")
	if err != nil {
		log.Fatal(err)
	}
	outJs, err := os.Create("./themes/" + theme.Name + ".ts")
	if err != nil {
		log.Fatal(err)
	}
	err = templateJs.Execute(outJs, params)
	if err != nil {
		log.Fatal(err)
	}

	outCss.Close()
}

type Style struct {
	Color     *string
	FontStyle *string
	Prio      *int
}

type TokenColorSettings struct {
	Foreground *string
	FontStyle  *string
}

type TokenColor struct {
	Scope    interface{}
	Settings TokenColorSettings
}

type VsCodeTheme struct {
	Colors      map[string]string
	TokenColors []TokenColor
}

func find(data VsCodeTheme, keys ...string) Style {
	style := Style{}

	for _, key := range keys {
		if value, exist := data.Colors[key]; exist {
			return Style{Color: &value}
		}

		for _, tokenColor := range data.TokenColors {
			scopes := []string{}
			rt := reflect.TypeOf(tokenColor.Scope)
			if tokenColor.Scope == nil {
				continue
			}

			switch rt.Kind() {
			case reflect.Slice:
				for _, s := range tokenColor.Scope.([]interface{}) {
					scopes = append(scopes, s.(string))
				}
			case reflect.String:
				splitted := strings.Split(tokenColor.Scope.(string), ",")
				for _, s := range splitted {
					scopes = append(scopes, strings.TrimSpace(s))
				}
			default:
				panic(fmt.Sprintf("Unecpected scope type %s", rt))
			}

			for i, scope := range scopes {
				if scope == key &&
					(style.Color == nil || *style.Prio > i) &&
					tokenColor.Settings.Foreground != nil {
					style.Color = tokenColor.Settings.Foreground
					style.Prio = &i
				}

				if scope == key && style.FontStyle == nil && tokenColor.Settings.FontStyle != nil {
					style.FontStyle = tokenColor.Settings.FontStyle
				}
			}
		}
	}

	if style.Color == nil {
		panic(fmt.Sprintf("Could not find color by: %s", keys))
	}

	return style
}

type TemplateParams struct {
	ExportPrefix string
	Dark         bool

	// Editor
	Background         Style
	Foreground         Style
	Selection          Style
	Cursor             Style
	DropdownBackground Style
	DropdownBorder     Style
	ActiveLine         Style
	MatchingBracket    Style

	// Syntax
	Keyword   Style // if else, etc
	Storage   Style // const, let, etc - Not supported in CM
	Parameter Style // fn(parmater)    - Not supported in CM
	Variable  Style
	Function  Style
	String    Style
	Constant  Style // ???
	Type      Style // x: MyType
	Class     Style // class MyClass
	Number    Style
	Comment   Style
	Heading   Style
	Invalid   Style
	Regexp    Style
}

func makeTemplateParams(theme Theme, content []byte) TemplateParams {
	var data VsCodeTheme
	err := jsonc.Unmarshal(content, &data)

	if err != nil {
		log.Fatal("JSON parse error: ", err)
	}

	params := TemplateParams{
		ExportPrefix: theme.Name,
		Dark:         theme.Dark,
		// Layout
		// ========================================================================
		Background:         find(data, "editor.background"),
		Foreground:         find(data, "foreground", "input.foreground"),
		Selection:          find(data, "editor.selectionBackground"),
		Cursor:             find(data, "editorCursor.foreground", "foreground"),
		DropdownBackground: find(data, "editor.background"),
		DropdownBorder:     find(data, "dropdown.border", "foreground"),
		ActiveLine:         find(data, "editor.lineHighlightBackground", "editor.selectionBackground"),
		MatchingBracket:    find(data, "editorBracketMatch.background", "editor.lineHighlightBackground", "editor.selectionBackground"),
		// Syntax
		// ========================================================================
		Keyword:   find(data, "keyword"),
		Storage:   find(data, "storage", "keyword"),
		Variable:  find(data, "variable.parameter", "variable.other", "variable.language", "variable", "foreground"),
		Parameter: find(data, "variable.parameter", "variable.other", "variable"),
		Function:  find(data, "support.function", "support", "entity.name.function", "entity.name"),
		String:    find(data, "string"),
		Constant:  find(data, "constant", "constant.character", "constant.keyword"),
		Type:      find(data, "support.type", "support", "entity.name.class"),
		Class:     find(data, "entity.name.class", "entity.name"),
		Number:    find(data, "constant.numeric", "constant"),
		Comment:   find(data, "comment"),
		Heading:   find(data, "markup.heading", "markup.heading.setext"),
		Invalid:   find(data, "invalid", "editorError.foreground", "errorForeground", "foreground", "input.foreground"),
		Regexp:    find(data, "string.regexp", "string"),
	}

	return params
}

func extractTheme(theme Theme) ([]byte, error) {
	r, err := zip.OpenReader("extensions/" + theme.Path)

	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	for _, f := range r.File {
		if f.Name != theme.File {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		content, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Fatal(err)
		}

		rc.Close()

		return content, nil
	}

	return nil, fmt.Errorf("Cound not find file %s in extension", theme.File)
}

func KebabToCamelCase(kebab string) (camelCase string) {
	isToUpper := false
	for _, runeValue := range kebab {
		if isToUpper {
			camelCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			if runeValue == '-' {
				isToUpper = true
			} else {
				camelCase += string(runeValue)
			}
		}
	}
	return
}
