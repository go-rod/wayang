# Overview

![](https://img.shields.io/github/v/tag/go-rod/wayang?sort=semver)

Table of Contents
=================

* [Overview](#overview)
* [Table of Contents](#table-of-contents)
* [Running a program](#running-a-program)
* [The Action JSON Object](#the-action-json-object)
* [Program Structure](#program-structure)
    * [Selectors](#selectors)
    * [Actions](#actions)
    * [Steps](#steps)
* [Examples](#examples)
    * [Navigate to a website](#navigate-to-a-website)



Rod is a High-level Devtools driver directly based on [DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/).
It's designed for web automation and scraping. Rod also tries to expose low-level interfaces to users, 
so that whenever a function is missing users can easily send control requests to the browser directly.

Wayang is a controller of Rod, making it possible to write programs with Rod in a more language neutral way.
While it does not cover all of Rod's API at the moment, we plan to make it feature complete.

# Running a program  

1. Build the CLI version for wayang. You scan use the make file provided in the [CLI](./cli) folder. 
 
2. There is currently no support for reading programs from STDIN. This is a feature that is planned to be added soon.
Instead, to link a script to run, you must provide the location of a JSON file. e.g. `--file="example.json"`

3. Provide other optional arguments. 
`--headless=[true|false]` will allow you to specify whether or not to run wayang in headless mode. 
With headless mode enabled, Chrome runs in the background and is not rendered. 
`--outputFile` can also be used to write the program output to a file. 

4. Read the documentation. The current JSON project is in alpha and not fully tested. 
You can still see examples in our [parser test file](./impl_test.go)

# The Action JSON Object

Actions are the bread and butter of wayang. 
They are what allow you to use wayang for web scraping, testing, and automation.

Actions are all required to have a key:value pair called `action`. This tells wayang what it should do.
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
When working with the `navigate` action, wayang needs to know what is the location of the link you want to navigate to. 
It is *required* to provide a `link` key so that wayang can navigate you to the page you want to go to.   
```json
{
  "action": "navigate",
  "link": "https://google.com/"
}
```

The action blurb above shows how we first tell wayang that we want to perform the `navigate` action, 
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
When executing an action, you can refer to the selector by that name by using the the string `$selector_name`.

### Actions

Sometimes, you may have repeatable actions that you want to use multiple times throughout your code. 
Alongside the `do` action, you can have your own simple "macro actions" you u can use around your code to do repeatable actions. 
Similarly to selectors, you define your custom actions in the root `actions` body.
Custom actions by themselves do not run, and the order in which they are defined do not matter. 
They are only run when they are called on in the steps tag, or when another action calls it which is running. 

Currently, there is no feature to provide parameters/arguments to actions. 

If you have a suggestion on how we can implement this while maintaining the simplicity of this language, 
feel free to create a new issue with your proposed solution. We are always looking on ways to improve wayang.  

### Steps

Coming to the final part of the program structure, we have the steps. 
The steps are an **array** of actions, that are ran in order of declaration in the program.
In the example above, the action with the name `action type` (which doesn't exist, and is only being used as an example)
is first executed. The other key value pairs in the body will be used as arguments when that action is executed.

# Examples

### Navigate to a website

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
