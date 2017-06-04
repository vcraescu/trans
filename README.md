## About Trans

Package Trans implements a simple translator similar to the one from Symfony2/3. 

## Installation

Install in the usual way:

    go get -u github.com/vcraescu/trans

## Usage

translations/messages.en.yml
```yaml
forms:
    label:
        name: Name
        email: Email
    description: This is just some description %foo%
    button:
        ok: OK
        cancel: Cancel

core:
    message:
        success: Successfully updated
        error: Update failed
    flash:
        success: Success flash
```

translations/messages.ro.yml
```yaml
forms:
    label:
        name: Nume
        email: Email
    description: Aceasta este doar o descriere %foo%
    button:
        ok: OK
        cancel: Anuleaza

core:
    message:
        success: Modificarea a fost salvata cu succes
        error: Actualizarea a esuat
```

```go
wd, err := os.Getwd()
if err != nil {
    log.Fatalln(err)
}

dirname := wd + "/translations"

t := trans.Translator{Locale: "en", FallbackLocale: "ro"}
t.Load(dirname)

fmt.Println(t.T("messages.forms.description", trans.Locale("en"), trans.Params{"%foo%": "this is foo"}))
fmt.Println(t.T("messages.forms.label.name"))
fmt.Println(t.T("messages.core.flash.success", trans.Locale("ro")))
```

Please refer to the [GoDoc API](https://godoc.org/github.com/vcraescu/trans) listing for a summary of the API. 
