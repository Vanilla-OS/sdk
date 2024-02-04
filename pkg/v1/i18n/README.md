# i18n

Localization support for the Vanilla OS SDK.

## Extract strings

To extract strings from the source code, we use the `xspreak` tool. Can be
installed with:

```sh
go install github.com/vorlif/xspreak@latest
```

To extract the POT, run:

```sh
xspreak -D path/to/source/ -p path/to/source/locale
```
