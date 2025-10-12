# eep!

`eep` is a flexible library for building slideshows with Ebitengine.

## Why?

While Google Slides and Microsoft PowerPoint are certainly powerful,
I often find myself wanting to draw diagrams programmatically.
Such diagrams are painful to integrate into either of these, since
they need to be rasterized and then inserted into the presentation.

Beyond static diagrams, I'm excited by the ability to write programs
that draw very dynamic diagrams and visualizations!

Taking this a step further, you can go absolutely nuts and [write
shaders](https://ebitengine.org/en/documents/shader.html) to make your
presentation the coolest in the world.

## Why not revealjs?

I have used [revealjs](https://revealjs.com) and I liked it.
But what I really want is a WYSIWYG experience for some parts.
Like, for example, the text styling UI in Google Slides is so much
nicer thatn writing code.
Similarly, the graphics UX in Google Slides is very, very nice for
one-off graphics.

While this library isn't there yet, and maybe it never will get there,
I would love to be able to add an edit mode that lets you augment your
presentation with a similar UX for rapid editing.

Do the hard (fun) parts in the code, do the easy (boring) parts in a
nice UX.

## Design

`eep` sticks to the core philosophy of Ebitengine, which I like call
"shut up and `Draw`."
That is, the core primitive is just drawing things to an `ebiten.Image`,
just like it is with Ebitengine.

In the same way, an `eep.Slide` is just an interface that can draw to
the screen, with an optional update function if it's stateful.
`eep.Slide` also composes very well.
You can just draw right on the slide!

`eep` provides some very basic slideshow functionality, along with nicer
APIs for certain primitives.
In particular:
- `graphics`: This is a vector graphics package heavily inspired by
  https://pkg.go.dev/github.com/fogleman/gg but built on top of Ebiten's
  vector graphics for efficient GPU rasterization.
- `text`: This is a package for drawing text to the screen, but at a
  higher level than what Ebitengine offers.
  It's API is meant to reflect the same model that you might expect
  from Google Slides or Microsoft PowerPoint: a text box that can
  contain text with a variety of styles.
- `font`: This is a package for managing font sources and font faces.
  It contains some conveniences, such as being able to load fonts
  that are available on your system (comparably to Microsoft PowerPoint).
  This functionality is provided through
  https://pkg.go.dev/github.com/adrg/sysfont.  

## Backward compatibility

`eep` is currently still experimental.
Don't expect API compatibility.

## Disclaimer

This project has the Go license, but this project is neither an official
Google product nor an official Go package.
I just work on the Go team and it's the easiest way to cover my work.
I do not work on this as part of my job, this is purely a hobby project.
Do not expect support.
I may respond to issues and PRs, but do not expect me to.
This is for fun.

## TODO

- Each line in a text box should be able to have its own alignment.
- More convenience functions in `graphics`.
- More convenience functions in `text`.
- Bold and italicized fonts are a bit awkward because we don't
  manage font families.
  You have to explicitly request the corresponding font by knowing
  its name.
- More themes.
- More slide templates.
- Speaker notes (in separate window).
- Presentation timer (in separate window).
- Edit mode for WYSIWYG editing, which is serialized back into an overlay
  file that's loaded on startup.
- Math TeX expression support, via codeberg.org/go-latex/latex.
  (Most of the work is hooking up Ebiten fonts to what latex expects.)
