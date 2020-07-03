# Wayang

![](https://img.shields.io/github/v/tag/go-rod/wayang?sort=semver)

Table of Contents
=================

* [Wayang](#wayang)
* [Table of Contents](#table-of-contents)
* [Overview](#overview)
    * [What is Wayang used for?](#what-is-wayang-used-for)
* [Running a program](#running-a-program)
* [Examples](#examples)
    * [Navigate to a website](#navigate-to-a-website)
    * [Execute a custom action](#execute-a-custom-action)
    * [Use a custom selector](#use-a-custom-selector)
    * [Execute multiple actions in place of one](#execute-multiple-actions-in-place-of-one)
    * [Using an if action](#using-an-if-action)
* [The Action JSON Object](#the-action-json-object)
* [Program Structure](#program-structure)
    * [Selectors](#selectors)
    * [Actions](#actions)
      * [Steps](#steps)
* [Documentation](#documentation)
  * [Selector elements](#selector-elements)
      * [PROPOSED CHANGES](#proposed-changes)
  * [Special Actions](#special-actions)
    * [do](#do)
    * [forEach](#foreach)
    * [if](#if)
    * [store](#store)
  * [Text Result Actions](#text-result-actions)
    * [attribute](#attribute)
    * [html](#html)
    * [text](#text)
  * [Boolean Result Actions](#boolean-result-actions)
    * [has](#has)
    * [not](#not)
    * [textContains](#textcontains)
    * [textEqual](#textequal)
    * [textNotEqual](#textnotequal)
    * [visible](#visible)
  * [Generic Actions](#generic-actions)
      * [Note](#note)
    * [blur](#blur)
    * [clear](#clear)
    * [click](#click)
    * [error](#error)
    * [eval](#eval)
    * [focus](#focus)
    * [input](#input)
    * [log](#log)
    * [logStore](#logstore)
    * [navigate](#navigate)
    * [press](#press)
    * [scrollIntoView](#scrollintoview)
    * [selectAll](#selectall)
  * [Sleep/Wait Actions](#sleepwait-actions)
      * [Note](#note-1)
    * [sleep](#sleep)
    * [waitIdle](#waitidle)
    * [waitInvisible](#waitinvisible)
    * [waitLoad](#waitload)
    * [waitStable](#waitstable)
    * [waitVisible](#waitvisible)

# Overview

Rod is a High-level Devtools driver directly based on [DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/).
It's designed for web automation and scraping. Rod also tries to expose low-level interfaces to users, 
so that whenever a function is missing users can easily send control requests to the browser directly.

Wayang is a controller of Rod, making it possible to write programs with Rod in a more language neutral way.
While it does not cover all of Rod's API at the moment, we plan to make it feature complete.

### What is Wayang used for?

Wayang was created for primarily integration testing. 
It allows QA engineers to write web automation tests more easily.
We also plan to add more features to make it easier for web automation developers to use this library. 

# Running a program  

1. Build the CLI version for Wayang. You can use the make file provided in the [CLI](cli) folder. 
 
2. There is currently no support for reading programs from STDIN. This is a feature that is planned to be added soon.
Instead, to link a script to run, you must provide the location of a JSON file. e.g. `--file="example.json"`

3. Provide other optional arguments. 
`--headless=[true|false]` will allow you to specify whether or not to run Wayang in headless mode. 
With headless mode enabled, Chrome runs in the background and is not rendered. 
`--outputFile` can also be used to write the program output to a file. 

4. Read the documentation. The current JSON project is in alpha and not fully tested. 
You can still see examples in our [parser test file](./impl_test.go)

# Examples

### Navigate to a website
A basic example showing how we can access a different website.

```json
{
  "steps": [
    {
      "action": "navigate",
      "link": "https://google.com"
    }
  ]
}
```

### Execute a custom action
We define a custom action that navigates to google. 
We then refer and run that action by prefixing the action name with a `$`.  

```json
{
  "actions": {
    "custom": {
      "action": "navigate",
      "link": "https://google.com"
    }
  },
  "steps": [
    {
      "action": "$custom"
    }
  ]
}
```

### Use a custom selector
This example shows defining a custom selector that can be reused throughout the program.
This is very useful when a website updates their selectors. You don't need to manually update every action. 

```json
{
  "selectors": {
    "gSearch": "//input[@name='q']"
  },
  "steps": [
    {
      "action": "input",
      "element": "$gSearch",
      "text": "wayang"
    }
  ]
}
```

### Execute multiple actions in place of one
In this example, we define a custom action, which normally only allow you to execute a single action. 
Using the `do` action we are able to chain multiple statements together.
 
```json
{
  "actions": {
    "custom": {
      "action": "do",
      "statements": [
        {
          "action": "input",
          "element": "//input[@type'text']",
          "text": "my text"
        },
        {
          "action": "blur",
          "element": "//input[@type='text']"
        }
      ]
    }
  },
  "steps": [
    {
      "action": "$custom"
    }
  ]
}
```

### Using an `if` action
Here we test our script completes it's actions by testing for error elements or checking for some text to be equal.
If you want to do more complex things, you can also use more complex xpath queries to test text, contains, and more.

```json
{
  "action": "if",
  "condition": {
    "action": "has",
    "element": "//div[@id='error']"
  },
  "execute": {
    "action": "error",
    "message": "expected error node to not be present"
  },
  "otherwise": {
    "action": "log",
    "message": "program successfully completed"
  }
}
```

# The Action JSON Object

Actions are the bread and butter of Wayang. 
They are what allow you to use Wayang for web scraping, testing, and automation.

Actions are all required to have a key:value pair called `action`. This tells Wayang what it should do.
A few examples of actions are:
- navigate
- click
- sleep
- input

The above actions also need more information though. 
When using the `navigate` action, you must provide a value for the `link` key.
When using the `click` action, you must provide a value for the `element` key. 
When using the `sleep` action, you must provide a value for the `duration` key (in seconds). 
For the `input` action, a `text` key is required, but a `element` key is optional. 
(More information available in the documentation for each action below)

Actions can also be used to test logic, and retrieve information. 
The `if` action depends on 2 other actions, and an optional third action.
It depends on a child `condition` action, which responds with either `true` or `false` based on the result of the action.
 
For example, the `has` action will return a true/false value depending on whether or not the page contains an element 
based on the element key:value pair. This can be used in conjunction with the `if` action to allow logic in your programs.
This is most useful when writing tests for a website, where you can log information and behaviour using `if` statements.

When working with actions, it is important to know what values are used to complete the action.
When working with the `navigate` action, Wayang needs to know what is the location of the link you want to navigate to. 
It is *required* to provide a `link` key so that Wayang can navigate you to the page you want to go to.   
```json
{
  "action": "navigate",
  "link": "https://google.com/"
}
```

The action blurb above shows how we first tell Wayang that we want to perform the `navigate` action, 
and then provide it information about where to navigate to.

Here is another example showing how the more complicated action `if` works.
The program requires us to provide the `condition` sub-action, as well as a `statement` sub-action.
If the `condition` action returns true, the `statement` action is run. *but*, if the action is not true,
we can provide an *optional* `otherwise` sub-action which will be run. If an `otherwise` sub-action is not provided,
the program will simply move on. 

```json
{
  "action": "if",
  "condition": {
      "action": "has",
      "element": "$myErrorNode"
  },
  "execute": {
      "action": "click",
      "element": "$errorNodeCloseButton"
  },
  "otherwise": {
    "action": "log",
    "message": "program successfully ran"
  }
}
```

The example above will check if the page has an element which matches the `myErrorNode` variable. 
If such a variable exists, we want to dismiss the error by clicking the element that matches the`errorNodeCloseButton`
element. We can also couple this with a `do` action (which allows you to run multiple actions in place of one), to perhaps
try our action again. The sky is the limit with what you can achieve. 

# Program Structure

The basic structure of a basic program is as follows:
```json
{
  "selectors": {
    "selector_name": "selector_value",
    "...": "..."
  },
  "actions": {
    "action_name": {
      "action": "action type",
      "other_arguments": "...",
      "...": "..."
    },
    "...": {}
  },
  "steps": [
    {
      "action": "action type",
      "element": "selector_name",
      "other_arguments": ""
    },
    {
      "action": "$action_name",
      "...": "..."
    }
  ]
}
```

There's a lot going on here, so let's go through it step by step.

### Selectors

Firstly, we define our custom selectors. 
This can be really handy when you want to refactor selectors in future code without going through your entire script.
***Note: Selectors must be defined using XPath notation.*** You can learn more about how XPath works [here](https://www.educba.com/xpath-operators/). 

Selectors is a map of key-pair values which are queried through the use of the dollar symbol. 
In the example above, we define a selector with the name `selector_name`. 
When executing an action, you can refer to the selector by that name by using the string `$selector_name`.

### Actions

Sometimes, you may have repeatable actions that you want to use multiple times throughout your code. 
Alongside the `do` action, you can have your own simple "macro actions" 
you can use around your code to do repeatable actions. 
Similarly to selectors, you define your custom actions in the root `actions` body.
Custom actions by themselves do not run, and the order in which they are defined do not matter. 
They are only run when they are called on in the steps tag, or when another action calls it which is running. 

Currently, there is no feature to provide parameters/arguments to actions. 

If you have a suggestion on how we can implement this while maintaining the simplicity of this language, 
feel free to create a new issue with your proposed solution. We are always looking on ways to improve Wayang.  

#### Steps

Coming to the final part of the program structure, we have the steps. 
The steps are an **array** of actions, that are ran in order of declaration in the program.
In the example above, the action with the name `action type` (which doesn't exist, and is only being used as an example)
is first executed. The other key value pairs in the body will be used as arguments when that action is executed.

# Documentation

## Selector elements

Currently, selector elements can either be an xpath string, or a reference to a defined selector in the 
`selectors` block. If the selector starts with a `$`, Wayang will parse it as a custom selector query.

#### PROPOSED CHANGES
We plan to add support for CSS Selectors soon. Here is our proposal to how we are planning to support them
in a backwards-compatible matter.

* If the value of a selector is type string
    * If it starts with a `$` symbol, parse it as a selector query, and use the result of its value.
    * Otherwise, parse it as an xpath query
* If the value of a selector is a block
    * Look for the two (required) sub-parameters:
        * `by`
            * If the `by` parameter is `xpath`, `x`, or `xp`, parse the selector as a xpath query
            * If the `by` parameter is `css` or `c`, parse the selector as a CSS Selector query.
        * `value`
            * The value of the selector that will be queried. Placeholders here are NOT allowed.

**Note**:

When defining global usage selectors, the above proposed changes will also be applied.

```json
[
  {
    "element": "this will be an xpath statement"
  },
  {
    "element": "$this will query the global selectors"
  },
  {
    "element": {
      "by": "xpath",
      "value": "this will be an xpath statement"
    }
  },
  {
    "element": {
      "by": "css",
      "value": "this will be a css query"
    }
  },
  {
    "element": {
      "by": "css",
      "value": "$this is not going to query global selectors"
    }
  }
]
```

## Special Actions

### do

Execute multiple actions in place of a single action. 
This can be useful when used in custom actions to execute multiple statements, and `if` actions.

**Parameters**:
- `statements` The list of actions that will be executed in order of declaration.
    - Required: Yes
    - Type: array(Action)

**Returns**: The result of the last statement. 

```json
{
  "action": "do",
  "statements": [
    {
      "action": "a"
    },
    {
      "action": "b"
    }
  ]
}
```

Action `a` will be ran first (as it is declared first), and then action `b` will be ran.
The result of action `b` will be the result of the `do` statement.

### forEach

WIP

### if

Execute an action conditionally.
This can be useful when writing tests to check for the presence of elements, or if they are visible.

**Parameters**:
- `condition`: An action, whose result will determine whether `statement` or `otherwise` is ran.
    - Required: Yes
    - Type: Action &rarr; bool
- `statement`: The action ran when `condition` is true.
    - Required: Yes
    - Type: Action
- `otherwise`: The action ran when `condition` is false.
    - Required: No
    - Type: Action 

**Returns**: The result of the action of `statement` or `otherwise` (Based on if `condition` is true or not). 
If the condition is false and `otherwise` is not provided, it will return `nil`.  

```json
{
  "action": "if",
  "condition": {
    "action": "has",
    "element": "//div[@id='error']"
  },
  "execute": {
    "action": "error",
    "message": "expected error node to not be present"
  },
  "otherwise": {
    "action": "log",
    "message": "program successfully completed"
  }
}
```

Wayang will check if `//div[@id='error]` is present on the website. If it is present, 
the program will error with the message `expected error node to not be present`. 
Otherwise, it will log `program successfully completed`, and continue/exit. 

### store

Store information into the program environment.

**Parameters**:
- `items`: A map of key-value pairs where the values will be added to the store. 
The values in the map do not have to be actions, and can also link to custom actions or statements.
    - Type: map[string] &rarr; Anything
    - Required: No. (Why are you executing this action then?)
    
**Returns**: `nil`

*NOTE*: This is not used as a place to retrieve and query information later yet. 
At the moment, this is only used as a way to get output information after running a program.
We are planning to add support to allow the querying of the store in places such as if statements. 

*NOTE*: Using store will allow you to overwrite previously stored items in the program.

```json
{
  "action": "store",
  "items": {
    "responseCode": {
      "action": "text",
      "element": "$responseCode"
    },
    "success": true
  }
}
```

Here we store two items into our program for output; the text of the `$responseCode` element, 
and `true` for the `success` key. 

## Text Result Actions

### attribute

The attribute action allows you to query the value of an attribute of an element. 
It will always be as the string representation.

```html
<p class="paragraph">...</p>
```

In the `p` tag above, the `class` attribute will return `paragraph`.

**Parameters**:
- `element`: The element that the attribute is queried from.
    - Type: selector
    - Required: Yes
- `name`: The name of the attribute to query.
    - Type: string
    - Required: Yes

**Returns**: The value of the queried attribute as a string.

```json
{
  "action": "attribute",
  "element": "//p",
  "name": "class"
}
``` 

### html

The html action returns the raw outer HTML representation of the element. 
It will also include all child elements present within the element.

**Parameters**:
- `element`: The specified element whose html is returned.
    - Type: selector
    - Required: Yes
    
**Returns**: The value of the outer HTML of the element as a string.

```json
{
  "action": "html",
  "element": "/html"
}
```

The above action will return the html of the entire document as a string (assuming it is a proper formatted website).

### text

This action will return the value for `input` and `textarea` elements, join the selected options for a `select` element,
and return the inner text for all other elements.

**Parameters**:
- `element`: The specified element the text is queried from.
    - Type: selector
    - Required: Yes
    
**Returns**: The text value of the element.

```json
{
  "action": "text",
  "element": "//input[@type='submit']"
}
```

## Boolean Result Actions

### has

The `has` action returns if the specified element exists or not on the page. 
This also returns true if the element is not visible. To check if an element is visible, use the `visible` action.

**Parameters**:
- `element`: The element that is queried. 
    - Type: selector
    - Required: Yes
    
**Returns**: A boolean value depending on the presence of the element on the page.

```json
{
  "action": "has",
  "element": "//div[@id='success-message']"
}
```

### not

The `not` action will result in the inverse of the provided action's result.
This also returns the inverse of a boolean.

**Parameters**:
- `statement`: The action whose inverse is returned.
    - Type: 
        - Action &rarr; bool
        - bool
    - Required: Yes

**Returns**: The inverse of the provided Action/Boolean.

```json
{
  "action": "not",
  "statement": {
    "action": "has",
    "element": "//div[contains(@class, 'success')]"
  }
}
```

### textContains

This action will return whether or not the result of the action contains the specified test.

**Parameters**:
- `statement`: The action whose result is tested to contain the text.
    - Type: Action &rarr; string
    - Required: Yes
- `text`: The text that is tested to be contained in the statement result.
    - Type: string
    - Required: Yes
- `ignoreCase`: Whether or not to ignore the case of `statement` and `text`.
    - Type: bool
    - Required: No
    - Default: `false`

**Returns**: Whether or not `statement`'s result contains the `text` substring.

```json
{
  "action": "textContains",
  "text": "welcome",
  "element": "//div[@class='intro-text']", 
  "ignoreCase": "false"
}
```

### textEqual

The `textEqual` action returns whether or not the result of the `statement` action equals the `text` parameter.

**Parameters**:
- `statement`: The action whose result is equal to the `text` parameter.
    - Type: Action &rarr; string
    - Required: Yes
- `text`: The text that is tested with the result of `statement`.
    - Type: string
    - Required: Yes
- `ignoreCase`: Whether or not to ignore the case of `statement` and `text`
    - Type: bool
    - Required: No
    - Default: `false`

**Returns**: Whether or not `statement`'s result equals to `text`.

```json
{
  "action": "textEqual",
  "statement": {
    "action": "attribute",
    "element": "//*[@id='result']",
    "name": "score"
  },
  "text": "100%"
}
```

### textNotEqual

The inverse result of `textEqual`

See [textEqual](#textequal)

```json
{
  "action": "textNotEqual",
  "statement": {
    "action": "attribute",
    "element": "//*[@id='result']",
    "name": "score"
  },
  "text": "100%"
}
```

### visible

The `visible` action returns whether or not the provided parameter `element` is visible on the page or not.

*NOTE*: If the element is not found on the page (visible or not), the program *will* wait indefinitely until it is present.
Support for a timeout is planned soon. 

Here is how the visibility of an element is determined (JavaScript):
```js
function visible() {
      const box = this.getBoundingClientRect()
      const style = window.getComputedStyle(this)
      return style.display !== 'none' &&
        style.visibility !== 'hidden' &&
        !!(box.top || box.bottom || box.width || box.height)
    }
```

**Parameters**:
- `element`: The element whose visibility is checked. 
    - Type: Action &rarr; string
    - Required: Yes

**Returns**: whether provided element is visible on the page.

```json
{
  "action": "visible",
  "element": "$errorNode"
}
```

## Generic Actions

#### Note

All the general actions defined below will always return nil, unless explicitly specified otherwise. 

### blur

Unfocus the element selector provided. This will also call the `onblur` functions defined for any element.

**Parameters**:
- `element`: The element to blur
    - Type: selector
    - Required: Yes        

```json
{
  "action": "blur",
  "element": "//input[@type='text']"
}
```

### clear

Clear will empty the text of an `input` or a `textarea`.
This can also be achieved using a `selectAll` action with a `input` with empty text.

**Parameters**:
- `element`: The element to clear the text from
    - Type: selector
    - Required: Yes

```json
{
  "action": "clear",
  "element": "//input[@type='text']"
}
```

### click

Click an element. This will also call the `click` event for most elements.

**Parameters**:
- `element`: The element to click
    - Type: selector
    - Required: Yes

```json
{
  "action": "click",
  "element": "//input@[type='submit']"
}
```

### error

Exit the program in an error state. This will also print the error message provided. 
The program will not execute any other actions after an error action executes.

**Parameters**:
- `message`: The reason or cause for the error
    - Type: string
    - Required: Yes

```json
{
  "action": "error",
  "message": "The website does not have a 'success' element."
}
```

### eval

Execute javascript code on the page. The javascript can also be executed in the context of another element. 

**Parameters**:
- `expression`": The code that will be executed. You must provide a function to be executed.
    - Type: string
    - Required: Yes
- `element`: A possible element to use for the context of the function. The element can be accessed with `this`.
    - Type: selector
    - Required: No

**Returns**: The raw value of the result of the expression. For string results, it wil also return the quotes. 

```json
[
  {
    "action": "eval",
    "expression": "() => console.log('hello!')"
  },
  {
    "action": "eval",
    "element": "//input[@type='submit']",
    "expression": "() => this.submit()"
  }
]
```

### focus

Focus an element. This will also call the `focus` event for the element.

**Parameters**:
- `element`: The element to focus
    - Type: selector
    - Required: Yes

```json
{
  "action": "focus",
  "element": "//input[@id='username']"
}
```

### input

Insert text into an input element. If an element isn't provided, 
it will instead press the keys into the page as you would with a keyboard.

**Parameters**:
- `element`: A possible element to input text into.
    - Type: selector
    - Required: No 
- `text`": The text that will be inputted into an element, or pressed.
    - Type: string
    - Required: Yes
    
```json
[
  {
    "action": "input",
    "element": "//input[@type='text']",
    "text": "Hello "
  },
  {
    "action": "input",
    "text": "World!"
  }
]
```

The second input action will still input the text into the `//input[@type='text']` element, as it is still focused.

### log

Print text to the console/logger output. 

**Parameters**:
- `message`: The message to output.
    - Type: string
    - Required: Yes

```json
{
  "action": "log",
  "message": "Successfully completed action!"
}
```

### logStore

Log an action to the console/logger from the program store.

**Parameters**:
- `key`: The key of the value to output from the store, prefixed with a $. 
    - Type: string
    - Required: Yes

```json
{
  "action": "logStore",
  "key": "$responseCode"
}
```

### navigate

Change the URL and load a new website. The request will block until the initial page response is complete. 
Not until everything loads.

**Parameters**:
- `link`: The website to connect to.
    - Type: string
    - Required: Yes
 
```json
{
  "action": "navigate",
  "link": "https://google.com"
  
}
```

### press

The `press` action will input any key into the element provided, or into the page (refer to the `input` action).
The press action gives you more control to specific keys, such as the `enter` key.

[Full list of supported keys](https://gist.github.com/Hamzantal/e4f465712caf0a444433db387b2f60a6)

**Parameters**:
- `element`: A possible element to input text into.
    - Type: selector
    - Required: No 
- `key`": The character/rune that will be inputted into an element, or pressed.
    - Type:
        - string (Refer to the gist above)
        - rune (e.g '\u0102')
    - Required: Yes

### scrollIntoView

Scroll the element into view, if it is not currently in the viewport.

**Parameters**:
- `element`: The element to scroll into view.
    - Type: selector
    - Required: Yes 

```json
{
  "action": "scrollIntoView",
  "element": "//div[contains(@class, 'footer')]"
}
```

### selectAll

Call the `select` method on a HTML element. On input elements, it selects the value. 

```js
function selectAllText () {
  this.select()
}
```

**Parameters**:
- `element`: The element to select all the text of.
    - Type: selector
    - Required: Yes 

```json
{
  "action": "selectAll",
  "element": "//input[@type='text']"
}
```

## Sleep/Wait Actions

#### Note

All the sleep/wait actions defined below will always return nil, unless explicitly specified otherwise. 

### sleep
### waitIdle
### waitInvisible
### waitLoad
### waitStable
### waitVisible
