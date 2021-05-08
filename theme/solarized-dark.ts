import {EditorView} from '@codemirror/view'
import {Extension} from '@codemirror/state'
import {HighlightStyle, tags as t} from '@codemirror/highlight'

export const config = {
  name: 'solarizedDark',
  dark: true,
  background: '#002b36',
  foreground: '#93a1a1',
  selection: '#073642',
  cursor: '#839496',
  dropdownBackground: '#00212b',
  dropdownBorder: '#2aa19899',
  keyword: '#859900',
  storage: '#93A1A1',
  variable: '#839496',
  parameter: '',
  function: '#268BD2',
  string: '#2AA198',
  constant: '#CB4B16',
  type: '#859900',
  class: '#CB4B16',
  number: '#D33682',
  comment: '#657B83',
  heading: '#268BD2',
  invalid: '',
  regexp: '#D30102',
}

export const solarizedDarkTheme = EditorView.theme({
  '&': {
    color: config.foreground,
    backgroundColor: config.background,
    '& ::selection': {backgroundColor: config.selection},
    caretColor: config.cursor,
  },

  '&.cm-focused .cm-cursor': {borderLeftColor: config.cursor},
  '&.cm-focused .cm-selectionBackground, .cm-selectionBackground': {backgroundColor: config.selection},

  '.cm-panels': {backgroundColor: config.dropdownBackground, color: config.foreground},
  '.cm-panels.cm-panels-top': {borderBottom: '2px solid black'},
  '.cm-panels.cm-panels-bottom': {borderTop: '2px solid black'},

  '.cm-searchMatch': {
    backgroundColor: config.dropdownBackground,
    outline: `1px solid ${config.dropdownBorder}`
  },
  '.cm-searchMatch.cm-searchMatch-selected': {
    backgroundColor: config.selection
  },

  '.cm-activeLine': {backgroundColor: config.selection},
  '.cm-selectionMatch': {backgroundColor: config.selection},

  '.cm-matchingBracket, .cm-nonmatchingBracket': {
    backgroundColor: config.invalid,
    outline: 'none'
  },
  '.cm-gutters': {
    backgroundColor: config.background,
    color: config.foreground,
    border: 'none'
  },
  '.cm-lineNumbers, .cm-gutterElement': {color: 'inherit'},

  '.cm-foldPlaceholder': {
    backgroundColor: 'transparent',
    border: 'none',
    color: config.foreground
  },

  '.cm-tooltip': {
    border: `1px solid ${config.dropdownBorder}`,
    backgroundColor: config.dropdownBackground
  },
  '.cm-tooltip-autocomplete': {
    '& > ul > li[aria-selected]': {
      backgroundColor: config.selection,
      color: config.foreground
    }
  }
}, {dark: config.dark})

export const solarizedDarkHighlightStyle = HighlightStyle.define([
  // const, let, function, if
  {tag: t.keyword, color: config.keyword},
  // document
  {tag: [t.name, t.deleted, t.character, t.macroName], color: config.variable},
  // getElementById
  {tag: [t.propertyName], color: config.function},
  // "string"
  {tag: [t.processingInstruction, t.string, t.inserted, t.special(t.string)], color: config.string},
  // render
  {tag: [t.function(t.variableName), t.labelName], color: config.function},
  // ???
  {tag: [t.color, t.constant(t.name), t.standard(t.name)], color: config.constant},
  // btn, count, fn render()
  {tag: [t.definition(t.name), t.separator], color: config.variable},
  {tag: [t.className], color: config.class},
  {tag: [t.number, t.changed, t.annotation, t.modifier, t.self, t.namespace], color: config.number},
  {tag: [t.typeName], color: config.type, fontStyle: config.type},
  {tag: [t.operator, t.operatorKeyword], color: config.keyword},
  {tag: [t.url, t.escape, t.regexp, t.link], color: config.regexp},
  {tag: [t.meta, t.comment], color: config.comment},
  {tag: t.strong, fontWeight: 'bold'},
  {tag: t.emphasis, fontStyle: 'italic'},
  {tag: t.link, textDecoration: 'underline'},
  {tag: t.heading, fontWeight: 'bold', color: config.heading},
  {tag: [t.atom, t.bool, t.special(t.variableName)], color: config.variable},
  {tag: t.invalid, color: config.invalid},
])

export const solarizedDark: Extension = [
  solarizedDarkTheme,
  solarizedDarkHighlightStyle,
]
