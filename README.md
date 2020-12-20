# webshot

A simple utility using chromedp and golang to screenshot a HTML file.

I realized that make some designs in HTML + CSS is nice but the problem is render them to a image.

For convenience it starts a webserver putting your file as `/` and exposes the current folder if the browser is not asking for `/`.
Chromedp, that is a remote control for Google Chrome take care to request `/` in this webserver then screenshot the page. Any dependency like `image.png` is handled as it is in the same folder.

# Dependencies

- Same as [chromedp](https://github.com/chromedp/chromedp)
