# XMLParser: dead simple & lightweight strict XML parser

## What is this?

This is a tiny XML parser  that parses basic XML files written in about 170 SLOC of Go

## Why?

1. Because everyone and their grandma seeks to write an XML parser these days
2. I wanted to write some Go

## What does it do?

It can parse tags, tag attributes, tag content (with some quirks, like no whitespace/NL stripping). It's a strict parser, so it will strictly check values of all attribute keys, of all tag names, and will ensure all opened tags are closed.

# What doesn't it do?

- Comments
- Anything not mentionned above and part of the XML spec
- Clean lexing. It instead relies on if/elses. Very primitive.

## How to use it?

See example.go and sample-data.yml

## Should I use it?

Probably not. Use encoding/xml instead.

## Where performance?

See encoding/xml. Also, performance on XML docs, really?