# Markdown Resume Builder
This tool was developed as a learning project centered around Markdown parsing and GUIs.
It streamlines the process of generating a resume from basic markdown, supplemented with tidbits of HTML.

![](resume-builder.gif)

## Features
- Generates a HTML or PDF resume styled using [GitHub CSS](style.css) from a markdown file
- Usable through both GUI or CLI
- Supports including a FontAwesome kit to include icons
- Supports including a picture

### Libraries used
- [Goldmark](https://github.com/yuin/goldmark) for Markdown parsing
- [ChromeDP](https://github.com/chromedp/chromedp) to use the Chrome browser to print to PDF in a headless fashion
- [Fyne](https://github.com/chromedp/chromedp) for the GUI
- [GitHub Markdown CSS](https://github.com/sindresorhus/github-markdown-css) for the styling

## Usage
### Prerequisites
- The [binaries](https://github.com/jlvoiseux/markdown-resume-builder/releases) were only generated and tested for Windows. 
- A Chrome install is required
### CLI
```
.\markdown-resume-builder.exe <SOURCE_RESUME_PATH> <html|pdf> <DESTINATION_FOLDER> <FONTAWESOME_KIT> <PICTURE_PATH>
```
The last two CLI arguments are optional.