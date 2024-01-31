# Go Template Manager

GTM is a Go library that enables you to work with the templates easily!

It handles layout-based templates and enables you to write your next Monolith application faster with Go!

Some of the features that GTM provides are:

- Layouts
- Includes
- Components
- Functions
- Vite

GTM uses Go's templates under the hood and you can write normal Go templates!

## Installation

To install GTM you need to run the following command:

```sh
go get github.com/saeedvaziry/gtm
```

## Init View

In your main app you need to initialize the views and load the templates.

```go
package main

import (
	"github.com/saeedvaziry/gtm"
)

func main() {
	gtm.Init(&gtm.View{
		BaseDir: "resources/views", // default
		ComponentsDir: "components", // default
		Extension: "html", // default
	})
	gtm.Load() // optional
}
```

`gtm.Load()` is optional to load the views. If you don't call this it will get the view from the files in the runtime when you want to render the view.

**BaseDir**: The base directory that you want to store your views and by default it is `resources/views`.

**ComponentsDir**: Inside the `BaseDir` you can asign a directory for your components and by default is `components` so it means your components should be stored in `resources/views/components`.

**Extension**: The file types you want to use as views and load them. By defult is `html`.

## Render

The `Render` function will render the view and returns the string so you can send it to the output!

```go
output := gtm.Render("home/index.html", nil)
```

### Render with data

```go
output := gtm.Render("home/index.html", map[string]interface{}{
    "title": "Go Template Manager"
})
```

Example:

```html
<div>
    {{ .title }}
</div>
```

### Render with Data and Functions

```go
output := gtm.Render("home/index.html", map[string]interface{}{
    "title": "Go Template Manager",
    "sayHi": func() string {
        return "Hi!"
    },
})
```

The function `sayHi` will be available in the view and you can use it by calling it inside the view file like `{{ sayHi }}`.

Example:

```html
<div>
    {{ sayHi }}
</div>
```

## Layouts

GTM is here to make the layouts as easy as possible! Here is an example of using the layouts:

```html
resources/views/layouts/base.html

<html>
    <head>
    </head>
    <body>
        <nav>...</nav>
        <header>...</header>
        <main>
            @child('content')
        </main>
        <footer></footer>
    </body>
</html>
```

```html
resources/views/home/index.html

@layout('layouts/base.html')

@section('content')
<h1>Hello World!</h1>
@end
```

## Includes (Injections)

Inside every view file you can include another one by using the `@include`

```html
resources/views/layouts/partials/footer.html

<footer>
    This is the footer!
</footer>
```

```html
resources/views/layouts/base.html

<html>
    <head>
    </head>
    <body>
        <nav>...</nav>
        <header>...</header>
        <main>
            @child('content')
        </main>
        @include('layouts/partials/footer.html')
    </body>
</html>
```

## Components

Components are powerful tools that can help you to reduce the code and reuse them!

To use the components you need to create them inside the `components` folder or where you defined the `ComponentsDir`.

GTM converts every component to a tag that starts with `x-`.

Here is an example:

```html
resources/views/components/primary-button.html

<button class="py-2 px-4 rounded-md bg-indigo-600 text-white text-center">
    @slot
</button>
```

And to use this component:

```html
resources/views/any-view-file.html

<x-primary-button>Click Me!</x-primary-button>
```

## Vite

GTM additionally supports Vite for easier development! All you need to do is to configure your resources with Vite and watch your resource files on a Vite Dev Server in your local machine. Then you can configure Vite with GTM as bellow:

```go
package main

import (
	"github.com/saeedvaziry/gtm"
)

func main() {
	gtm.ViteInit(&gtm.Vite{
		BuildPath: "static/build", // default
		DevServerUrl: "http://localhost:3000", // default
	})
	gtm.AddFunc("asset", gtm.ViteAsset)
}
```

This will add the `asset` function into all of your view files when they are being rendered.

And in your view files you can import the assets. For example:

```html
<link rel="stylesheet" href='{{ asset "resources/css/app.css" }}' />
```

If Vite's Dev Server was running on the URL you configured (`DevServerUrl`) then it will import the asset from there otherwise it will read the `manifest.json` file in your static path that Vite generates and will get the correct path of the asset!

For example you can see this output when the Vite Dev Server is up:

```html
<link rel="stylesheet" href='http://localhost:3000/static/build/resources/css/app.css' />
```

And this when you build and ship to production:

```html
<link rel="stylesheet" href='/static/build/css/app.3247289.css' />
```
