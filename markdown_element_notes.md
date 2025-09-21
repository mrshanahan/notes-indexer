Types of Markdown elements:
- Whitespace:
  - For formatting purposes, it appears that tabs are treated as 4 spaces - that is, one tab indentation will create a new nested list, but two tab indentations will append the current line to the previous list item (see "Lists")
  - More specifically, tabs used as indentation are converted to 4 spaces as a preprocessing step. This is evidenced by the behavior when creating a new code block via whitespace after a list: if the list is unordered & has no leading whitespace before the '-', and the code block is created following it with two tabs, then the resulting code block will have two leading spaces (2 spaces for the list + 4 spaces for the code block + 2 space prefix = 8 spaces) despite the fact that no spaces were used in indenting the code block
- Paragraphs:
  - <= 3 leading whitespace
  - Paragraphs are separated by blank lines - two paragraphs placed one after the other are is effectively one line
- Lists:
  - Unordered lists:
    - Can use '*' or '-'
    - Can mix if necessary
  - Ordered lists: 
    - Starting glyph matches `\d+\.`
    - First item in list determines starting order - e.g.:

          1. foo    ->    1. foo
          1. bar    ->    2. bar

          2. foo    ->    2. foo
          2. bar    ->    3. bar

  - Indentation behavior seems to vary, but it appears to follow from where the first non-space characters appear in the parent list item. Thus:
    - First level: establishes standard for second level
    - Second level: 1-3 spaces from beginning of first-level non-space characters; establishes standard for third level
    - Third level: 1-3 spaces from beginning of second-level non-space characters; establishes standard for fourth level
    - etc.
    - Anything beyond 3 appears to append it to the previous list item instead of establishing a next-level item
      - NB: If then after a 5+ space prepended item the next line is a 2-5 space indented item, it behaves as if it is the next item. E.g.:

            - Test1
                  - Test2
              - Test3

            ->

            * Test1 - Test2
                * Test3

  - Must be followed by blank line (or else following line is just appended to last one)
    - That is, the end-of-list marker is "\n\n", whereas "\n(?!\n)" means a line continuation for the previous bullet
  - Blank line separating list items creates separate lists
  - Sub-sections of Markdown can start treating the leftmost non-start-glyph character as the leftmost margin - works for:
    - Preformatted sections
    - Tables
- Headers:
  - Preceded by '#' symbol; any amount of space (so long as it's not a preformatted section) allowed
  - '#' may be specified to whatever depth - h5?
- In-paragraph formatting:
  - `*` or `_`: italics
  - `**` or `__`: bold
  - `/*(**)+/` or `/_(__)+/`: italics & bold
  - `/**(**)+/` or `/__(__)+/`: bold
  - `~~`: strikethrough
  - ``/`+/``: preformatted
  - `[<text>](<url>)`: link
    - Plain URLs are auto-formatted as links
  - Rules:
    - If one of the glyphs is followed by a non-space character, and then another one of the same is followed by another non-space character, that segment receives the formatting
    - The search for the next one is non-greedy - in the following case, there are two separate italics segments, the first one covering "foo bar" and the second covering "boo baz":

          *foo bar* bing bang *boo baz* bom

    - Formatted sections can be nested, but cannot be overlapping without nesting (just like HTML). For example, the following section will result in "foo bar bing bang" being in italics and "bing" also being in bold:

          *foo bar **bing** bang* boom

      And the following would just have "foo bar bing bang" in italics, with the literal "**" glyphs rendered without bold formatting on the text in between:

          *foo bar **bing bang* boom baz** bom

    - Glyphs act as a beginning or ending marker depending on their position relative to the nearest word:
      - One that is preceded by a non-word character &amp; followed by a word character (i.e. at a *left* word boundary) is a **start** character, e.g. `This is a *start formatting glyph ...`
      - One that is preceded by a word character &amp; follwed by a non-word character (i.e. at a *right* word boundary) is an **end** character, e.g. `... and this is an end* formatting glyph`
      - One that is preceded &amp; follwed by word characters (i.e. *within* a word boundary) can act as **either a start or an end** character, e.g. `This glyph aff*ects the mi*ddle of a phrase`.
    - Note that the "word boundary" that the glyphs **approximately** respect Unicode-aware rules of the `\w` regex character class when determining where word boundaries are. "Approximately", because there are a few caveats:
      - Underscore (`_`) is treated as a word character, so the `_` inline formatting item cannot be used within a word. (This isn't really a problem as `_` is just an alternative `*` and `*` still works just fine for this purpose.)
      - Certain characters, like `'` and `,`, appear to trigger different behavior in different parsers - some treat these as word characters, and some do not.
        - **I don't see any reason why we can't just use whitespace &amp; line breaks as the boundary - there's really no reason to include punctuation.**
- Preformatting
  - Preformatted sections are defined to be any section that is separated by blank lines on top & bottom and are at least four spaces from the current indentation level
  - Four spaces are counted from the leftmost margin for text sections & from the beginning of the previous list glyph for list sections
  - Preformatted sections can also be designated with a "```" glyph with less than 4 leading spaces from the margin
- Block quotes
  - Anything beginning with ">" (preceded by <= 3 blank spaces)
  - All standard Markdown is valid inside a blockquote level - e.g. four spaces from ">" produces a 
  - Additional ">" characters produce additional block quote levels
  - Each subsequent line 
- Tables: TBD
- Footnotes: TBD

Paragraph

		Test

# - Is this a list?
#Is this a header?
# We keep *processing* inline elements!
    # Is this a header?

1. Buh - 1. 2. Hi
1. 
1. Test

- # Test
- Test1
  - # Test2
  Test3
    - Test3
    - Test4
    - # Test5
      - Test6

- > Test
  > test2
  - Test3

- 
- bloch


- Test
      - Test2
  - Test3

- Test1
  - Test2

    ```
    bloch
    ```
    black

- Test1
    - Test2

>     test

> test

- List
- List

1. a
1. b
1. c

1. test


1. lsdkfjsd

5. lkdjfkd

- test

1. lkdfjkd
2. lksdfjklsd

```
test
```

[link]

[link](https://link.com)


One paragraph
- Test

-Test test test

-- thing


This **is**a testab

This i__s a t__est.

This __is a__ test

This isn*_t a t*est.

This isn*'t a t*est.

This isn'*t a t*est.

Thi*s isn'*t a test.