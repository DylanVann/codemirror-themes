import { EditorView } from '@codemirror/view'
import type { Extension } from '@codemirror/state'
import { HighlightStyle, tags as t } from '@codemirror/highlight'

export const config = {
  name: 'base',
  dark: false,
  background: 'var(--cm-background)',
  foreground: 'var(--cm-foreground)',
  selection: 'var(--cm-selection)',
  cursor: 'var(--cm-cursor)',
  dropdownBackground: 'var(--cm-dropdown-background)',
  dropdownBorder: 'var(--cm-dropdown-border)',
  activeLine: 'var(--cm-active-line)',
  matchingBracket: 'var(--cm-matching-bracket)',
  keyword: 'var(--cm-keyword)',
  storage: 'var(--cm-storage)',
  variable: 'var(--cm-variable)',
  parameter: 'var(--cm-parameter)',
  function: 'var(--cm-function)',
  string: 'var(--cm-string)',
  constant: 'var(--cm-constant)',
  type: 'var(--cm-type)',
  class: 'var(--cm-class)',
  number: 'var(--cm-number)',
  comment: 'var(--cm-comment)',
  heading: 'var(--cm-heading)',
  invalid: 'var(--cm-invalid)',
  regexp: 'var(--cm-regexp)',
}

export const theme = EditorView.theme(
  {
    '&': {
      color: config.foreground,
      backgroundColor: config.background,
    },

    '.cm-content': { caretColor: config.cursor },

    '&.cm-focused .cm-cursor': { borderLeftColor: config.cursor },
    '&.cm-focused .cm-selectionBackground, .cm-selectionBackground, & ::selection':
      { backgroundColor: config.selection },

    '.cm-panels': {
      backgroundColor: config.dropdownBackground,
      color: config.foreground,
    },
    '.cm-panels.cm-panels-top': { borderBottom: '2px solid black' },
    '.cm-panels.cm-panels-bottom': { borderTop: '2px solid black' },

    '.cm-searchMatch': {
      backgroundColor: config.dropdownBackground,
      outline: `1px solid ${config.dropdownBorder}`,
    },
    '.cm-searchMatch.cm-searchMatch-selected': {
      backgroundColor: config.selection,
    },

    '.cm-activeLine': { backgroundColor: config.activeLine },
    '.cm-activeLineGutter': { backgroundColor: config.background },
    '.cm-selectionMatch': { backgroundColor: config.selection },

    '.cm-matchingBracket, .cm-nonmatchingBracket': {
      backgroundColor: config.matchingBracket,
      outline: 'none',
    },
    '.cm-gutters': {
      backgroundColor: config.background,
      color: config.foreground,
      border: 'none',
    },
    '.cm-lineNumbers, .cm-gutterElement': { color: 'inherit' },

    '.cm-foldPlaceholder': {
      backgroundColor: 'transparent',
      border: 'none',
      color: config.foreground,
    },

    '.cm-tooltip': {
      border: `1px solid ${config.dropdownBorder}`,
      backgroundColor: config.dropdownBackground,
      color: config.foreground,
    },
    '.cm-tooltip.cm-tooltip-autocomplete': {
      '& > ul > li[aria-selected]': {
        background: config.selection,
        color: config.foreground,
      },
    },
  },
  { dark: config.dark },
)

export const highlightStyle = HighlightStyle.define([
  // const, let, function, if
  { tag: t.keyword, color: config.keyword },
  // document
  {
    tag: [t.name, t.deleted, t.character, t.macroName],
    color: config.variable,
  },
  // getElementById
  { tag: [t.propertyName], color: config.function },
  // "string"
  {
    tag: [t.processingInstruction, t.string, t.inserted, t.special(t.string)],
    color: config.string,
  },
  // render
  { tag: [t.function(t.variableName), t.labelName], color: config.function },
  // ???
  {
    tag: [t.color, t.constant(t.name), t.standard(t.name)],
    color: config.constant,
  },
  // btn, count, fn render()
  { tag: [t.definition(t.name), t.separator], color: config.variable },
  { tag: [t.className], color: config.class },
  {
    tag: [t.number, t.changed, t.annotation, t.modifier, t.self, t.namespace],
    color: config.number,
  },
  { tag: [t.typeName], color: config.type, fontStyle: config.type },
  { tag: [t.operator, t.operatorKeyword], color: config.keyword },
  { tag: [t.url, t.escape, t.regexp, t.link], color: config.regexp },
  { tag: [t.meta, t.comment], color: config.comment },
  { tag: t.strong, fontWeight: 'bold' },
  { tag: t.emphasis, fontStyle: 'italic' },
  { tag: t.link, textDecoration: 'underline' },
  { tag: t.heading, fontWeight: 'bold', color: config.heading },
  { tag: [t.atom, t.bool, t.special(t.variableName)], color: config.variable },
  { tag: t.invalid, color: config.invalid },
  { tag: t.strikethrough, textDecoration: 'line-through' },
])

export const extension: Extension = [theme, highlightStyle]
