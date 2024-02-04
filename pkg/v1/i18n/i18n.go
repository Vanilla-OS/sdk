package i18n

import (
	"fmt"
	"strings"

	appTypes "github.com/vanilla-os/sdk/pkg/v1/app/types"

	"github.com/vorlif/spreak"
	"golang.org/x/text/language"
)

// NewLocalizer creates a new localizer for the application, the localizer is
// used to localize strings in the application. The localizer requires the
// application and the locale to be passed in, the locale is the language
// the user wants to use. If the locale is not found or is empty, the
// localizer will default to English, assuming it as the fallback.
//
// Example:
//
//	t, err := i18n.NewLocalizer(app, "")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Println(t.Get("I am Batman!"))
func NewLocalizer(app *appTypes.App, locale string) (*spreak.Localizer, error) {
	foundLocale, err := language.Parse(locale)
	if err != nil {
		foundLocale = language.English
	}

	// we need to get the supported languages from the locales file system
	// to do so we expect a LINGUAS file to be present
	linguas, err := app.LocalesFS.ReadFile("locales/LINGUAS")
	if err != nil {
		return nil, fmt.Errorf("no LINGUAS file found: %v", err)
	}

	// spreak.WithLanguage requires a slice of interfaces
	supportedLanguages := make([]interface{}, 0)
	for _, l := range strings.Split(string(linguas), "\n") {
		if l == "" {
			continue
		}
		supportedLanguages = append(supportedLanguages, language.MustParse(l))
	}

	// we need to create a new bundle for the localizer, here we use the RDNN
	// as the default localizer domain
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain(app.RDNN),
		spreak.WithDomainFs(app.RDNN, app.LocalesFS),
		spreak.WithRequiredLanguage(foundLocale),
		spreak.WithLanguage(supportedLanguages...),
	)
	if err != nil {
		return nil, err
	}

	return spreak.NewLocalizer(bundle, foundLocale), nil
}
