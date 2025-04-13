# Wishlister

Wishlister is a small webapp that enable a user to create and share wishlists.
The goal is to have a small, free, open-source, bloat-free, without affiliation link of wathever, simple solution.

Wishlister try to make the process as easy as possible.
One consequence of this is that there no account creation.
A user can provides an email address, but that is not mandatory.

## Features

Currently, it is possible to create a list, edit it (with an admin link that must not be shared), and view it
(with a link that can be shared).

On creation, the user can provide an email address which is used to send them the two links.

## Roadmap

I want to add some more features:
* add a way to ask to be sent again all the list created with a given email address: the list would be sent to
  the given email address, there should be no security issue (assuming that the inital creator has still control
  over their email address)
* rework the UI, maybe with Beer CSS: I don't really like the look of Pico CSS, especially on desktop

The next features could be implemented if there is a willing for them, but as long as I am the only one deploying
this app, it will probably never be the case:
* make all the hardcoded values (email server, website name's, etc) configurable
* add internationalization: the app UI is currently in french only

The next features will never be implemented:
* account creation: the goal of this app is to be as simple as possible, creating an account is not simple
* any suggestion of existing products or affiliation links to e-commerce: this is exactly what I dislike with
  existing wishlist apps

## Live deployment

This app is deployed at https://www.malistdevoeux.fr, in french.

## Publish a new package version

This is mostly a memo for myself:
```
git tag vX.Y
git push origin tag vX.Y
VERSION=X.Y make publish
