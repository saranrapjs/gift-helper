# gift-helper

This is a gift registry app (e.g. for a baby or wedding registry) that uses Google Forms to record gift purchases.

## Google Form setup

This assumes you have a Google Form that's wired up to a Google Spreadsheet somewhere — [here's a guide that explains how to do that](https://support.google.com/docs/answer/2917686?hl=en).

This app assumes that your Google form has three text fields, in this very specific order:

- name
- notes
- URL

The app will fetch the first three field names from the Google form at startup, and submit responses to those three fields in that order. If you _dont_ have those three fields things may break. Naming the fields is unimportant, but they should be text fields.

## What it does

On startup, the server fetches the fields from the Google form. It will also process the template file (optionally override-able via the `GIFT_TEMPLATEDIR` directory, where it will look for a `template.html`) and use any Go template tags with the format `{{.Bought "https://purchaseable-item-url/"}}` to determine which things can be bought. This tag will be replaced with `data-bought="boolean" href="https://purchaseable-item-url/"` at render time, which the Javascript that's included in the example will use to change how the item is displayed.

There are a few other environment variables that get pulled into the template, made available as a convenience (`GIFT_ADDRESS`, `GIFT_NAME`, `GIFT_EMAIL`) but you may find it easier to just edit these by hand in the HTML. The template is loaded afresh with every page load, so the process doesn't need to restart for changes to the template to appear.

Responses will be persisted to a local `tracking.json` file (so that this information survives across restarts of the process), as well as sent along to the Google Form.

## usage

The only required environment variable is `GIFT_GOOGLEFORM`, which points to the Google form that responses will be recorded in. All other config options live in the `Config` struct.

```bash
GIFT_GOOGLEFORM="your-google-form-url-here" go run main.go
```
