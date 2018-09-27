Date: 2017-10-11
Title: Mobile viewport recipe
cat: dev

If your mobile css is not sizing divs properly, try these two quick things:

###Add the viewport meta to the <head> of your HTML:

```
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
```

###Use percentages for div size:

```
<div style="display:inline-block;width:45%;">
```

![mobile-screen](/images/mobile-screen.png)
