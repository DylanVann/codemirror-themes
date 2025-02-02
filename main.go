package main

import (
  "fmt"
  "log"
  "os"
  "io"
  "strings"
  "reflect"
  "text/template"
  "net/http"
  "archive/zip"
  "io/ioutil"
  "muzzammil.xyz/jsonc"
)

type Theme struct {
  Name string
  Url string
  File string
  Dark bool
}

func main() {
  themes := []Theme{
    {
      Name: "dracula",
      Url: "https://marketplace.visualstudio.com/_apis/public/gallery/publishers/dracula-theme/vsextensions/theme-dracula/2.22.3/vspackage",
      File: "extension/theme/dracula.json",
      Dark: true,
    },
    {
      Name: "solarized-light",
      Url: "https://marketplace.visualstudio.com/_apis/public/gallery/publishers/ryanolsonx/vsextensions/solarized/2.0.3/vspackage",
      File: "extension/themes/light-color-theme.json",
      Dark: false,
    },
    {
      Name: "solarized-dark",
      Url: "https://marketplace.visualstudio.com/_apis/public/gallery/publishers/ryanolsonx/vsextensions/solarized/2.0.3/vspackage",
      File: "extension/themes/dark-color-theme.json",
      Dark: true,
    },
    {
      Name: "material-light",
      Url: "https://marketplace.visualstudio.com/_apis/public/gallery/publishers/Equinusocio/vsextensions/vsc-material-theme/33.1.2/vspackage",
      File: "extension/build/themes/Material-Theme-Lighter.json",
      Dark: false,
    },
    {
      Name: "material-dark",
      Url: "https://marketplace.visualstudio.com/_apis/public/gallery/publishers/Equinusocio/vsextensions/vsc-material-theme/33.1.2/vspackage",
      File: "extension/build/themes/Material-Theme-Default.json",
      Dark: true,
    },
    {
      Name: "github-light",
      Url: "https://marketplace.visualstudio.com/_apis/public/gallery/publishers/GitHub/vsextensions/github-vscode-theme/3.0.0/vspackage",
      File: "extension/themes/light.json",
      Dark: false,
    },
    {
      Name: "github-dark",
      Url: "https://marketplace.visualstudio.com/_apis/public/gallery/publishers/GitHub/vsextensions/github-vscode-theme/3.0.0/vspackage",
      File: "extension/themes/dark.json",
      Dark: true,
    },
  }

  for _, theme := range themes {
    fmt.Println("Process theme: ", theme.Name)
    if _, err := os.Stat("./tmp/" + theme.Name + ".zip"); os.IsNotExist(err) {
      fmt.Println("  Download theme")
      downloadTheme(theme)
    }

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

  t, err := template.ParseFiles("./template.js")
  if err != nil {
    log.Fatal(err)
  }

  out, err := os.Create("./theme/" + theme.Name + ".ts")
  if err != nil {
    log.Fatal(err)
  }

  err = t.Execute(out, params)
  if err != nil {
    log.Fatal(err)
  }

  out.Close()
}

type Style struct{
  Color *string
  FontStyle *string
  Prio *int
}

type TokenColorSettings struct{
  Foreground *string
  FontStyle *string
}

type TokenColor struct{
  Scope interface{}
  Settings TokenColorSettings
}

type VsCodeTheme struct{
  Colors map[string]string
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
        if (
          scope == key &&
          (style.Color == nil || *style.Prio > i) &&
          tokenColor.Settings.Foreground != nil) {
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
  ExportPrefix       string
  Dark               bool

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
  Keyword            Style // if else, etc
  Storage            Style // const, let, etc - Not supported in CM
  Parameter          Style // fn(parmater)    - Not supported in CM
  Variable           Style
  Function           Style
  String             Style
  Constant           Style // ???
  Type               Style // x: MyType
  Class              Style // class MyClass
  Number             Style
  Comment            Style
  Heading            Style
  Invalid            Style
  Regexp             Style
}

func makeTemplateParams(theme Theme, content []byte) TemplateParams {
  var data VsCodeTheme
  err := jsonc.Unmarshal(content, &data)

  if err != nil {
    log.Fatal("JSON parse error: ", err)
  }

  params := TemplateParams{
    ExportPrefix:       KebabToCamelCase(theme.Name),
    Dark:               theme.Dark,
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
    Keyword:            find(data, "keyword"),
    Storage:            find(data, "storage", "keyword"),
    Variable:           find(data, "variable.parameter", "variable.other", "variable.language", "variable", "foreground"),
    Parameter:          find(data, "variable.parameter", "variable.other", "variable"),
    Function:           find(data, "support.function", "support", "entity.name.function", "entity.name"),
    String:             find(data, "string"),
    Constant:           find(data, "constant", "constant.character", "constant.keyword"),
    Type:               find(data, "support.type", "support", "entity.name.class"),
    Class:              find(data, "entity.name.class", "entity.name"),
    Number:             find(data, "constant.numeric", "constant"),
    Comment:            find(data, "comment"),
    Heading:            find(data, "markup.heading", "markup.heading.setext"),
    Invalid:            find(data, "invalid", "editorError.foreground", "errorForeground", "foreground", "input.foreground"),
    Regexp:             find(data, "string.regexp", "string"),
  }

  return params
}

func extractTheme(theme Theme) ([]byte, error) {
  r, err := zip.OpenReader("tmp/" + theme.Name + ".zip")

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

func downloadTheme(theme Theme) {
  resp, err := http.Get(theme.Url)
  if err != nil {
    log.Fatal(err)
  }

  defer resp.Body.Close()
  if resp.StatusCode != 200 {
    log.Fatal("Could not download theme: ", theme, "StatusCode: ", resp.StatusCode)
  }

  _ = os.Mkdir("./tmp", 0700)

  // Create the file
  out, err := os.Create("tmp/" + theme.Name + ".zip")
  if err != nil {
    log.Fatal(err)
  }

  defer out.Close()

  _, err = io.Copy(out, resp.Body)

  if err != nil {
    log.Fatal(err)
  }
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
